package aop

import (
	"reflect"

	"github.com/gogap/errors"
)

type Proxy struct {
	beanID string
	funcs  map[string]MethodMetadata
}

func NewProxy(beanID string) *Proxy {
	return &Proxy{
		beanID: beanID,
		funcs:  make(map[string]MethodMetadata),
	}
}

func (p *Proxy) BeanID() string {
	return p.beanID
}

func (p *Proxy) Method(fn interface{}) (method interface{}) {

	methodName := ""
	if methodMetadata, err := getMethodMetadata(fn); err != nil {
		panic(err)
	} else {
		methodName = methodMetadata.MethodName
	}

	if metadata, exist := p.funcs[methodName]; exist {
		method = metadata.method
		return
	}
	return
}

func (p *Proxy) Invoke(method interface{}, args ...interface{}) (result *InvokeResult) {

	methodName := ""
	if methodMetadata, err := getMethodMetadata(method); err != nil {
		result = &InvokeResult{
			beanID:     p.beanID,
			methodName: methodName,
			err:        err,
		}
	} else {
		methodName = methodMetadata.MethodName
	}

	fnMetadata, exist := p.funcs[methodName]
	if !exist {
		result = &InvokeResult{
			beanID:     p.beanID,
			methodName: methodName,
			err:        ErrInvokeFuncNotExist.New(errors.Params{"id": p.beanID, "method": methodName}),
		}
		return
	}

	fn := fnMetadata.method

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

func (p *Proxy) registryFunc(metadata MethodMetadata) {
	p.funcs[metadata.MethodName] = metadata
}
