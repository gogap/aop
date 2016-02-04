package aop

import (
	"github.com/rs/xid"
	"reflect"
)

var (
	errorType = reflect.TypeOf((*error)(nil)).Elem()
)

type AOP struct {
	aspects     []*Aspect
	beanFactory BeanFactory
}

func NewAOP() *AOP {
	return &AOP{}
}

func (p *AOP) SetBeanFactory(factory BeanFactory) {
	p.beanFactory = factory
}

func (p *AOP) AddAspect(aspect *Aspect) *AOP {
	p.aspects = append(p.aspects, aspect)
	return p
}

func (p *AOP) funcWrapper(bean *Bean, methodName string, methodType reflect.Type) func([]reflect.Value) []reflect.Value {
	beanValue := reflect.ValueOf(bean.instance)

	return func(inputs []reflect.Value) (ret []reflect.Value) {
		callID := xid.New().String()

		var err error
		defer func() {
			if err != nil {
				panic(err)
			}
		}()

		var args Args

		for _, arg := range inputs {
			args = append(args, arg.Interface())
		}

		joinPoint := JoinPoint{
			callID: callID,
			args:   args,
			target: bean,
		}

		errOutIndex := -1
		outLen := methodType.NumOut()

		if outLen > 0 {
			for i := 0; i < outLen; i++ {
				if methodType.Out(i) == errorType {
					errOutIndex = i
				}
			}
		}

		ret = make([]reflect.Value, outLen)
		for i := 0; i < outLen; i++ {
			ret[i] = reflect.Zero(methodType.Out(i))
		}

		var advicesGroup []map[AdviceOrdering][]*Advice
		for _, aspect := range p.aspects {
			var advices map[AdviceOrdering][]*Advice
			if advices, err = aspect.GetMatchedAdvices(bean, methodName, args); err != nil {
				return
			}

			advicesGroup = append(advicesGroup, advices)
		}

		callAdvicesFunc := func(order AdviceOrdering, resultValues ...reflect.Value) (e error) {

			for _, advices := range advicesGroup {
				if e = invokeAdvices(&joinPoint, advices[order], methodName, resultValues); e != nil {
					if errOutIndex >= 0 {
						ret[errOutIndex] = reflect.ValueOf(&e).Elem()
					}
					err = e
					return
				}
			}
			return
		}

		//@Before
		if err = callAdvicesFunc(Before); err != nil {
			return
		}

		//@Real func
		var retValues []reflect.Value

		funcInSturctName := getFuncNameByStructFuncName(methodName)

		realFunc := func(args ...interface{}) Result {
			values := []reflect.Value{}
			for _, arg := range args {
				values = append(values, reflect.ValueOf(arg))
			}

			return beanValue.MethodByName(funcInSturctName).Call(values)
		}

		//@Around
		var aroundAdvice *Advice
		for _, advices := range advicesGroup {
			if aroundAdvices, exist := advices[Around]; exist && len(aroundAdvices) > 0 {
				aroundAdvice = aroundAdvices[0]
				break
			}
		}

		if aroundAdvice != nil {
			pjp := ProceedingJoinPoint{JoinPointer: &joinPoint, method: realFunc}

			if err = invokeAdvices(&pjp, []*Advice{aroundAdvice}, methodName, nil); err != nil {
				if errOutIndex >= 0 {
					ret[errOutIndex] = reflect.ValueOf(&err).Elem()
				}
				return
			}
		} else {
			retValues = realFunc(inputs)
		}

		if IsTracing() {
			var metadata MethodMetadata
			if metadata, err = bean.methodMetadata(funcInSturctName); err != nil {
				return
			}

			appendTraceItem(callID, metadata.File, metadata.Line, "*"+funcInSturctName, methodName, bean.id)
		}

		defer func() {
			//@AfterPanic
			if e := recover(); e != nil {
				callAdvicesFunc(AfterPanic)
			}

			//@After
			callAdvicesFunc(After)
		}()

		if err != nil {
			//@AfterError
			callAdvicesFunc(AfterError)
		} else {
			//@AfterReturning
			callAdvicesFunc(AfterReturning, retValues...)
		}

		return retValues
	}
}

func (p *AOP) GetProxy(beanID string) (proxy *Proxy, err error) {
	var bean *Bean

	if bean, err = p.beanFactory.GetBean(beanID); err != nil {
		return
	}

	tmpProxy := NewProxy(beanID)

	beanValue := reflect.ValueOf(bean.instance)
	beanType := reflect.TypeOf(bean.instance)
	for i := 0; i < beanValue.NumMethod(); i++ {
		methodV := beanValue.Method(i)
		methodT := beanType.Method(i)

		mType := methodV.Type()

		var metadata MethodMetadata
		if metadata, err = getMethodMetadata(methodT); err != nil {
			return
		}

		newFunc := p.funcWrapper(bean, metadata.Method.Name, mType)
		funcV := reflect.MakeFunc(mType, newFunc)

		metadata.Method.Func = funcV // rewrite to new proxy func

		tmpProxy.registryFunc(metadata)
	}

	proxy = tmpProxy

	return
}
