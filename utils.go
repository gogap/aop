package aop

import (
	"github.com/gogap/errors"
	"reflect"
	"runtime"
	"strings"
)

func getMethodMetadata(method interface{}) (metadata MethodMetadata, err error) {
	if method == nil {
		err = ErrMethodIsNil.New()
		return
	}

	if reflect.TypeOf(method).Kind() != reflect.Func {
		err = ErrBadMethodType.New()
		return
	}

	v := reflect.ValueOf(method)

	pc := runtime.FuncForPC(v.Pointer())

	metadata.MethodName = strings.TrimRight(pc.Name(), "-fm")
	metadata.MethodName = strings.Replace(metadata.MethodName, "*", "", 1)
	metadata.File, metadata.Line = pc.FileLine(v.Pointer())
	metadata.method = method

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

		useJPArgs := false
		var jpArgs Args

		adviceArgsType := getFuncArgsType(advice.beanRef, advice.Method)
		lenAdviceArgs := len(adviceArgsType)

		if lenAdviceArgs != len(joinPoint.Args()) {
			useJPArgs = true
			jpArgs = make(Args, lenAdviceArgs)
		} else {
			for i, adviceArgType := range adviceArgsType {
				if adviceArgType != reflect.TypeOf(joinPoint.Args()[i]) {
					useJPArgs = true
					jpArgs = make([]interface{}, lenAdviceArgs)
				}
			}
		}

		// inject jp to advice args
		var invokeArgs Args
		if useJPArgs {
			jpType := reflect.TypeOf(joinPoint)
			retType := reflect.TypeOf(result)

			for i, argType := range adviceArgsType {
				if jpType.ConvertibleTo(argType) {
					jpArgs[i] = joinPoint
				} else if argType == retType {
					if advice.Ordering == AfterReturning {
						jpArgs[i] = result
					} else {
						panic(ErrJoinPointArgsUsage.New())
					}
				} else {
					panic(ErrUnknownJoinPointArgType.New(errors.Params{"id": advice.beanRef.id, "method": advice.Method, "refID": advice.PointcutRefID}))
				}
			}
			invokeArgs = jpArgs
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
