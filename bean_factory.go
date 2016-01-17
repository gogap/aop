package aop

import (
	"github.com/gogap/errors"
	"reflect"
)

type BeanFactory interface {
	RegisterBean(id string, value interface{}) BeanFactory
	GetBean(id string) (bean interface{}, err error)
}

type BeanFactoryAware interface {
	SetBeanFactory(factory BeanFactory)
}

type ClassicBeanFactory struct {
	beans map[string]interface{}
}

func NewClassicBeanFactory() BeanFactory {
	return &ClassicBeanFactory{
		beans: make(map[string]interface{}),
	}
}

func (p *ClassicBeanFactory) GetBean(id string) (bean interface{}, err error) {
	if id == "" {
		err = ErrBeanIDShouldNotBeEmpty.New()
		return
	}

	if v, exist := p.beans[id]; exist {
		return v, nil
	}

	return nil, ErrBeanNotExist.New(errors.Params{"id": id})
}

func (p *ClassicBeanFactory) RegisterBean(id string, beanObj interface{}) (factory BeanFactory) {
	var err error

	defer func() {
		if err != nil {
			panic(err)
		}
	}()

	if id == "" {
		err = ErrBeanIDShouldNotBeEmpty.New()
		return
	}

	if _, exist := p.beans[id]; exist {
		err = ErrBeanAlreadyRegistered.New(errors.Params{"id": id})
		return
	}

	if beanObj == nil {
		err = ErrBeanInstanceIsNil.New(errors.Params{"id": id})
		return
	}

	v := reflect.ValueOf(beanObj)
	if v.Kind() != reflect.Ptr {
		err = ErrBeanIsNotAnPtr.New(errors.Params{"id": id})
		return
	}

	p.beans[id] = beanObj

	return p
}
