package aop

import (
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

func (p *AOP) invokeAdvices(ordering AdviceOrdering, bean *Bean, methodName string, args Args) (err error) {
	var callAdvices []*Advice
	for _, aspect := range p.aspects {
		var advices []*Advice
		if advices, err = aspect.GetMatchedAdvices(ordering, bean, methodName, args); err != nil {
			return
		}

		callAdvices = append(callAdvices, advices...)
	}

	for _, advice := range callAdvices {
		var retFunc func()
		if _, err = advice.beanRef.Invoke(advice.Method, args, func(values ...interface{}) {
			if values != nil {
				for _, v := range values {
					if errV, ok := v.(error); ok {
						err = errV
					}
				}
			}

			if err != nil {
				return
			}
		}); err != nil {
			return
		}

		if retFunc != nil {
			retFunc()
		}
	}

	return
}

func (p *AOP) Invoke(beanID string, methodName string, args Args, callback ...interface{}) (err error) {
	var bean *Bean

	if bean, err = p.beanFactory.GetBean(beanID); err != nil {
		return
	}

	//@Before
	if err = p.invokeAdvices(Before, bean, methodName, args); err != nil {
		return
	}

	// Call Bean Service
	var retFunc func()
	retFunc, err = bean.MustInvoke(methodName, args, callback...)

	defer func() {
		//@AfterPanic
		if e := recover(); e != nil {
			p.invokeAdvices(AfterPanic, bean, methodName, args)
		}

		//@After
		err = p.invokeAdvices(After, bean, methodName, args)

		if err == nil && retFunc != nil {
			retFunc()
		}
	}()

	if err != nil {
		//@AfterError
		p.invokeAdvices(AfterError, bean, methodName, args)
	} else {
		//@AfterReturning
		p.invokeAdvices(AfterReturning, bean, methodName, args)
	}

	return
}

func (p *AOP) funcWrapper(bean *Bean, methodName string, methodType reflect.Type) func([]reflect.Value) []reflect.Value {
	beanValue := reflect.ValueOf(bean.instance)

	return func(inputs []reflect.Value) (ret []reflect.Value) {

		var args Args
		var err error
		for _, arg := range inputs {
			args = append(args, arg.Interface())
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

		//@Before
		if err = p.invokeAdvices(Before, bean, methodName, args); err != nil {
			if errOutIndex >= 0 {
				ret[errOutIndex] = reflect.ValueOf(&err).Elem()
			}
			return
		}

		retValues := beanValue.MethodByName(methodName).Call(inputs)

		defer func() {
			//@AfterPanic
			if e := recover(); e != nil {
				p.invokeAdvices(AfterPanic, bean, methodName, args)
			}

			//@After
			p.invokeAdvices(After, bean, methodName, args)
		}()

		if err != nil {
			//@AfterError
			p.invokeAdvices(AfterError, bean, methodName, args)
		} else {
			//@AfterReturning
			p.invokeAdvices(AfterReturning, bean, methodName, args)
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

	beanType := reflect.TypeOf(bean.instance)
	for i := 0; i < beanType.NumMethod(); i++ {
		method := beanType.Method(i)
		commonMethodType := getCommonFuncType(method)
		newFunc := p.funcWrapper(bean, method.Name, commonMethodType)
		funcV := reflect.MakeFunc(commonMethodType, newFunc)

		tmpProxy.registryFunc(method.Name, funcV.Interface())
	}

	proxy = tmpProxy

	return
}

// TODO: add type build cache
func getCommonFuncType(method reflect.Method) reflect.Type {
	inTypes := []reflect.Type{}
	for i := 0; i < method.Type.NumIn(); i++ {
		inTypes = append(inTypes, method.Type.In(i))
	}

	outTypes := []reflect.Type{}
	for i := 0; i < method.Type.NumOut(); i++ {
		outTypes = append(outTypes, method.Type.Out(i))
	}

	if len(inTypes) > 0 {
		inTypes = inTypes[1:]
	}

	return reflect.FuncOf(inTypes, outTypes, method.Type.IsVariadic())
}
