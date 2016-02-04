package aop

import (
	"reflect"
	"runtime"
	"strings"
)

func getMethodMetadata(method reflect.Method) (metadata MethodMetadata, err error) {
	if method.Func.IsNil() {
		err = ErrMethodIsNil.New()
		return
	}

	iMethod := method.Func.Interface()

	if reflect.TypeOf(iMethod).Kind() != reflect.Func {
		err = ErrBadMethodType.New()
		return
	}

	v := reflect.ValueOf(iMethod)

	pc := runtime.FuncForPC(v.Pointer())

	name := strings.TrimRight(pc.Name(), "-fm")
	name = strings.Replace(name, "*", "", 1)

	metadata.Method = method
	metadata.Method.Name = name
	metadata.File, metadata.Line = pc.FileLine(v.Pointer())

	return
}

func getFuncMetadata(fn interface{}) (metadata MethodMetadata, err error) {
	if fn == nil {
		err = ErrMethodIsNil.New()
		return
	}

	if reflect.TypeOf(fn).Kind() != reflect.Func {
		err = ErrBadMethodType.New()
		return
	}

	v := reflect.ValueOf(fn)

	pc := runtime.FuncForPC(v.Pointer())

	name := strings.TrimRight(pc.Name(), "-fm")
	name = strings.Replace(name, "*", "", 1)

	m := reflect.Method{
		Name:    name,
		PkgPath: "",
		Type:    reflect.TypeOf(fn),
		Func:    v,
		Index:   -1,
	}

	metadata.File, metadata.Line = pc.FileLine(v.Pointer())
	metadata.Method = m

	return
}

func invokeAdvices(joinPoint JoinPointer, advices []*Advice, methodName string, result Result) (err error) {
	for _, advice := range advices {
		var retFunc func()
		if IsTracing() {
			var metadata MethodMetadata
			if metadata, err = advice.beanRef.methodMetadata(advice.Method); err != nil {
				return
			}

			appendTraceItem(joinPoint.CallID(), metadata.File, metadata.Line, advice.Method, methodName, advice.beanRef.ID())
		}

		var adviceMethodMeta MethodMetadata
		if adviceMethodMeta, err = advice.beanRef.methodMetadata(advice.Method); err != nil {
			return
		}

		var invokeArgs Args
		if adviceMethodMeta.IsEqual(joinPointFuncType) ||
			adviceMethodMeta.IsEqual(proceedingJoinPointType) {
			invokeArgs = Args{joinPoint}
		} else if adviceMethodMeta.IsEqual(joinPointWithResultFuncType) {
			invokeArgs = Args{joinPoint, result}
		} else {
			invokeArgs = joinPoint.Args()
		}

		if _, err = advice.beanRef.Invoke(advice.Method, invokeArgs, func(values ...interface{}) {
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

func getFuncArgsType(bean *Bean, methodName string) (types []reflect.Type) {
	method := reflect.ValueOf(bean.instance).MethodByName(methodName)
	for i := 0; i < method.Type().NumIn(); i++ {
		types = append(types, method.Type().In(i))
	}

	return
}

func getFuncNameByStructFuncName(name string) string {
	if name == "" {
		return ""
	}

	index := strings.LastIndex(name, ".")
	if index > 0 {
		return name[index+1:]
	}
	return ""
}

func getFullStructName(v interface{}) (name string, err error) {
	beanT := reflect.TypeOf(v)
	name = beanT.String()
	name = strings.Replace(name, "*", "", 1)
	return
}
