package aop

import (
	"github.com/gogap/errors"
)

type BeanFactory interface {
	RegisterBean(id string, class string, value interface{}) BeanFactory
	GetBean(id string) (bean *Bean, err error)
}

type BeanFactoryAware interface {
	SetBeanFactory(factory BeanFactory)
}

type ClassicBeanFactory struct {
	beans map[string]*Bean
}

func NewClassicBeanFactory() BeanFactory {
	return &ClassicBeanFactory{
		beans: make(map[string]*Bean),
	}
}

func (p *ClassicBeanFactory) GetBean(id string) (bean *Bean, err error) {
	if id == "" {
		err = ErrBeanIDShouldNotBeEmpty.New()
		return
	}

	if v, exist := p.beans[id]; exist {
		return v, nil
	}

	return nil, ErrBeanNotExist.New(errors.Params{"id": id})
}

func (p *ClassicBeanFactory) RegisterBean(id string, class string, beanInstance interface{}) (factory BeanFactory) {
	var err error

	defer func() {
		if err != nil {
			panic(err)
		}
	}()

	if _, exist := p.beans[id]; exist {
		err = ErrBeanAlreadyRegistered.New(errors.Params{"id": id})
		return
	}

	var bean *Bean
	if bean, err = NewBean(id, class, beanInstance); err != nil {
		return
	}

	p.beans[id] = bean

	return p
}
