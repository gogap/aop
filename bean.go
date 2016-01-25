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

func (p *Bean) Invoke(methodName string, args Args, callback ...interface{}) (returnFunc func(), err error) {

	var beanType reflect.Type
	var beanValue reflect.Value

	var isSameArgs bool

	beanType = reflect.TypeOf(p.instance)
	beanValue = reflect.ValueOf(p.instance)

	if method, exist := beanType.MethodByName(methodName); !exist {
		err = ErrBeanMethodNotExit.New(errors.Params{"id": p.id, "class": p.class, "method": methodName})
		panic(err)
	} else {
		if method.Type.NumIn() == len(args)+1 {
			isSameArgs = true
			for i, arg := range args {
				tArg := reflect.TypeOf(arg)
				if tArg.String() != method.Func.Type().In(i+1).String() {
					isSameArgs = false
					break
				}
			}
		}

		compareArgCount := 0
		if reflect.Indirect(beanValue).Kind() == reflect.Struct {
			compareArgCount = compareArgCount + 1
		}

		if !isSameArgs && method.Type.NumIn() != compareArgCount {
			err = ErrWrongAdviceFuncArgsNum.New(errors.Params{"id": p.id, "class": p.class, "method": methodName})
			return
		}
	}

	var values []reflect.Value

	if isSameArgs {
		inputs := make([]reflect.Value, len(args))
		for i, _ := range args {
			inputs[i] = reflect.ValueOf(args[i])
		}
		values = beanValue.MethodByName(methodName).Call(inputs)
	} else {
		values = beanValue.MethodByName(methodName).Call([]reflect.Value{})
	}

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

func (p *Bean) MustInvoke(methodName string, args Args, callback ...interface{}) (returnFunc func(), err error) {
	var beanValue reflect.Value

	beanValue = reflect.ValueOf(p.instance)

	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
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
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}

	return beanValue.MethodByName(methodName).Call(inputs)
}
