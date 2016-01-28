package aop

import (
	"reflect"

	"github.com/gogap/errors"
)

type Bean struct {
	id       string
	class    string
	instance interface{}
}

func NewBean(id string, instance interface{}) (bean *Bean, err error) {
	if id == "" {
		err = ErrBeanIDShouldNotBeEmpty.New()
		return
	}

	if instance == nil {
		err = ErrBeanInstanceIsNil.New(errors.Params{"id": id})
		return
	}

	v := reflect.ValueOf(instance)
	if v.Kind() != reflect.Ptr {
		err = ErrBeanIsNotAnPtr.New(errors.Params{"id": id})
		return
	}

	class := ""
	if class, err = getFullStructName(instance); err != nil {
		return
	}

	bean = &Bean{
		id:       id,
		class:    class,
		instance: instance,
	}

	return
}

func (p *Bean) ID() string {
	return p.id
}

func (p *Bean) Class() string {
	return p.class
}

func (p *Bean) methodMetadata(methodName string) (metadata MethodMetadata, err error) {
	beanType := reflect.TypeOf(p.instance)

	var method reflect.Method
	exist := false
	if method, exist = beanType.MethodByName(methodName); !exist {
		err = ErrBeanMethodNotExit.New(errors.Params{"id": p.id, "class": p.class, "method": methodName})
		return
	}

	metadata, err = getMethodMetadata(method.Func.Interface())

	return
}

func (p *Bean) Invoke(methodName string, args Args, callback ...interface{}) (returnFunc func(), err error) {
	var beanValue reflect.Value

	beanValue = reflect.ValueOf(p.instance)

	inputs := make([]reflect.Value, len(args))

	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}

	values := beanValue.MethodByName(methodName).Call(inputs)

	if values != nil && len(values) > 0 {
		lastV := values[len(values)-1]
		if lastV.Interface() != nil {
			if errV, ok := lastV.Interface().(error); ok {
				if errV != nil {
					err = errV
					return
				}
			}
		}
	}

	if callback != nil && len(callback) > 0 {
		returnFunc = func() {
			reflect.ValueOf(callback[0]).Call(values)
		}
	}

	return
}

func (p *Bean) Call(methodName string, args Args) []reflect.Value {
	var beanValue reflect.Value

	beanValue = reflect.ValueOf(p.instance)

	inputs := make([]reflect.Value, len(args))
	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}

	return beanValue.MethodByName(methodName).Call(inputs)
}
