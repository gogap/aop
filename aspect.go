package aop

import (
	"github.com/gogap/errors"
)

type Args []interface{}

type Aspect struct {
	id          string
	beanRefID   string
	beanFactory BeanFactory

	pointcutIDs []string
	pointcuts   map[string]*Pointcut
	advices     []*Advice
}

func NewAspect(id, beanRefID string) *Aspect {
	if id == "" {
		panic(ErrAspectIDShouldNotBeEmpty.New())
	}

	if beanRefID == "" {
		panic(ErrBeanIDShouldNotBeEmpty.New())
	}

	return &Aspect{
		id:        id,
		beanRefID: beanRefID,
		pointcuts: make(map[string]*Pointcut),
	}
}

func (p *Aspect) ID() string {
	return p.id
}

func (p *Aspect) BeanRefID() string {
	return p.beanRefID
}

func (p *Aspect) AddPointcut(pointcut *Pointcut) *Aspect {
	if _, exist := p.pointcuts[pointcut.ID()]; !exist {
		p.pointcuts[pointcut.ID()] = pointcut
		p.pointcutIDs = append(p.pointcutIDs, pointcut.ID())
	}
	return p
}

func (p *Aspect) AddAdvice(advice *Advice) *Aspect {
	var beanRef *Bean
	var err error

	if beanRef, err = p.beanFactory.GetBean(p.beanRefID); err != nil {
		panic(err)
	}

	if advice.PointcutRefID == "" {
		err = ErrEmptyPointcutRefID.New()
		panic(err)
	}

	if pointcut, exist := p.pointcuts[advice.PointcutRefID]; exist {
		advice.pointcut = pointcut
	} else {
		panic(ErrPointcutNotExist.New(errors.Params{"id": advice.PointcutRefID}))
	}

	advice.beanRef = beanRef
	p.advices = append(p.advices, advice)

	return p
}

func (p *Aspect) SetBeanFactory(factory BeanFactory) {
	p.beanFactory = factory
	return
}

func (p *Aspect) GetMatchedAdvices(bean *Bean, methodName string, args Args) (advices map[AdviceOrdering][]*Advice, err error) {
	var advs map[AdviceOrdering][]*Advice = make(map[AdviceOrdering][]*Advice)

	for _, advice := range p.advices {
		if advice.pointcut.IsMatch(bean, methodName, args) {
			advs[advice.Ordering] = append(advs[advice.Ordering], advice)
		}
	}

	advices = advs

	return
}
