package aop

import (
	"reflect"

	"github.com/gogap/errors"
)

type Proxy struct {
	beanID string
	funcs  map[string]interface{}
}

func NewProxy(beanID string) *Proxy {
	return &Proxy{
		beanID: beanID,
		funcs:  make(map[string]interface{}),
	}
}

func (p *Proxy) BeanID() string {
	return p.beanID
}

func (p *Proxy) Method(name string) interface{} {
	fn, _ := p.funcs[name]
	return fn
}

func (p *Proxy) Invoke(methodName string, args ...interface{}) (result *InvokeResult) {

	fn, exist := p.funcs[methodName]
	if !exist {
		result = &InvokeResult{
			beanID:     p.beanID,
			methodName: methodName,
			err:        ErrInvokeFuncNotExist.New(errors.Params{"id": p.beanID, "method": methodName}),
		}
		return
	}

	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		result = &InvokeResult{
			beanID:     p.beanID,
			methodName: methodName,
			err:        ErrInvokeFuncTypeError.New(errors.Params{"id": p.beanID, "method": methodName}),
		}
		return
	}

	if fnType.NumIn() != len(args) {
		result = &InvokeResult{
			beanID:     p.beanID,
			methodName: methodName,
			err:        ErrWrongInvokeFuncArgsNum.New(errors.Params{"id": p.beanID, "method": methodName}),
		}
		return
	}

	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}

	fnValue := reflect.ValueOf(fn)

	result = &InvokeResult{
		beanID:     p.beanID,
		methodName: methodName,
		err:        nil,
		values:     fnValue.Call(inputs),
	}

	return
}

func (p *Proxy) registryFunc(name string, fn interface{}) {
	p.funcs[name] = fn
}
