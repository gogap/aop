package aop

import (
	"reflect"
	"sync"

	"github.com/gogap/errors"
)

type InvokeResult struct {
	beanID     string
	methodName string
	values     []reflect.Value
	err        error
	callOnce   sync.Once
	called     bool
}

func (p *InvokeResult) End(callback ...interface{}) (err error) {

	if p.called {
		return ErrEndInvokeTwice.New(errors.Params{"id": p.beanID, "method": p.methodName})
	}

	if p.err != nil {
		return p.err
	}

	p.callOnce.Do(func() {

		p.called = true

		if callback == nil || len(callback) == 0 {
			return
		}

		cbType := reflect.TypeOf(callback[0])
		if cbType.Kind() != reflect.Func {
			panic(ErrEndInvokeParamsIsNotFunc.New(errors.Params{"id": p.beanID, "method": p.methodName}))
		}

		if cbType.NumIn() != len(p.values) {
			panic(ErrWrongEndInvokeFuncArgsNum.New(errors.Params{"id": p.beanID, "method": p.methodName}))
		}

		cbValue := reflect.ValueOf(callback[0])

		cbValue.Call(p.values)
	})

	return
}

func (p *InvokeResult) MethodName() string {
	return p.methodName
}

func (p *InvokeResult) BeanID() string {
	return p.beanID
}

func (p *InvokeResult) Error() error {
	return p.err
}
