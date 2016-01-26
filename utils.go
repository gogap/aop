package aop

import (
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

func invokeAdvices(invokeID string, advices []*Advice, bean *Bean, methodName string, args Args) (err error) {
	for _, advice := range advices {
		var retFunc func()
		if IsTracing() {
			var metadata MethodMetadata
			if metadata, err = advice.beanRef.methodMetadata(advice.Method); err != nil {
				return
			}

			appendTraceItem(invokeID, metadata.File, metadata.Line, advice.Method, methodName, advice.beanRef.ID())
		}
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
