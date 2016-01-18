package aop

type Args []interface{}

type Aspect struct {
	id          string
	advices     map[AdviceOrdering][]*Advice
	beanRefID   string
	beanFactory BeanFactory
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
		advices:   make(map[AdviceOrdering][]*Advice),
	}
}

func (p *Aspect) ID() string {
	return p.id
}

func (p *Aspect) BeanRefID() string {
	return p.beanRefID
}

func (p *Aspect) AddAdvice(advice *Advice) *Aspect {
	var beanRef *Bean
	var err error

	if beanRef, err = p.beanFactory.GetBean(p.beanRefID); err != nil {
		panic(err)
	}

	advice.beanRef = beanRef
	p.advices[advice.Ordering] = append(p.advices[advice.Ordering], advice)
	return p
}

func (p *Aspect) SetBeanFactory(factory BeanFactory) {
	p.beanFactory = factory
	return
}

func filterAdvices(advices []*Advice) (matchedAdvices []*Advice, err error) {
	// check call stack, make sure not have cycle call

	return
}

func (p *Aspect) GetMatchedAdvices(ordering AdviceOrdering, bean *Bean, methodName string, args Args) (advices []*Advice, err error) {
	var advs []*Advice
	var exist bool

	if advs, exist = p.advices[ordering]; !exist {
		return
	}

	var retAdvs []*Advice

	for _, adv := range advs {
		var match bool
		if match, err = adv.IsMatch(ordering, bean, methodName, args); err != nil {
			return
		} else if match {
			retAdvs = append(retAdvs, adv)
		}
	}

	advices = retAdvs

	return
}
