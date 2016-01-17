package aop

import (
	"reflect"
)

type Args []interface{}

type Aspect struct {
	id          string
	advices     map[AdviceOrdering][]*Advice
	beanFactory BeanFactory
}

func NewAspect(id string) *Aspect {
	return &Aspect{
		id:      id,
		advices: make(map[AdviceOrdering][]*Advice),
	}
}

func (p *Aspect) SetBeanFactory(factory BeanFactory) {
	p.beanFactory = factory
}

func (p *Aspect) AddAdvice(advice Advice) *Aspect {
	return p
}

func (p *Aspect) Invoke(beanID string, methodName string, args Args, callback ...interface{}) (err error) {
	var bean interface{}
	var beanValue reflect.Value

	if bean, err = p.beanFactory.GetBean(beanID); err != nil {
		return
	}

	//TODO:

	//@Before

	//@After

	//@AfterReturning

	//@AfterError

	//@AfterPanic

	//@Around

	beanValue = reflect.ValueOf(bean)

	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}

	values := beanValue.MethodByName(methodName).Call(inputs)
	if callback != nil && len(callback) > 0 {
		reflect.ValueOf(callback[0]).Call(values)
	}

	return
}

func (p *Advice) getMatchedAdvices(ordering AdviceOrdering, beanID string, methodName string, args Args) (advices []*Advice, err error) {
	return
}
