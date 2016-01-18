package aop

type AOP struct {
	aspects     []*Aspect
	beanFactory BeanFactory
}

func NewAOP() *AOP {
	return &AOP{}
}

func (p *AOP) SetBeanFactory(factory BeanFactory) {
	p.beanFactory = factory
}

func (p *AOP) AddAspect(aspect *Aspect) *AOP {
	p.aspects = append(p.aspects, aspect)
	return p
}

func (p *AOP) invokeAdvices(ordering AdviceOrdering, bean *Bean, methodName string, args Args) (err error) {
	var callAdvices []*Advice
	for _, aspect := range p.aspects {
		var advices []*Advice
		if advices, err = aspect.GetMatchedAdvices(ordering, bean, methodName, args); err != nil {
			return
		}

		callAdvices = append(callAdvices, advices...)
	}

	for _, advice := range callAdvices {
		var retFunc func()
		if _, err = advice.beanRef.Invoke(advice.Method, args, func(values ...interface{}) {
			if values != nil {
				for _, v := range values {
					if errV, ok := v.(error); ok {
						err = errV
					}
				}
			}

			if err != nil {
				return
			}
		}); err != nil {
			return
		}

		if retFunc != nil {
			retFunc()
		}
	}

	return
}

func (p *AOP) Invoke(beanID string, methodName string, args Args, callback ...interface{}) (err error) {
	var bean *Bean

	if bean, err = p.beanFactory.GetBean(beanID); err != nil {
		return
	}

	//@Before
	if err = p.invokeAdvices(Before, bean, methodName, args); err != nil {
		return
	}

	// Call Bean Service
	var retFunc func()
	retFunc, err = bean.MustInvoke(methodName, args, callback...)

	defer func() {
		//@AfterPanic
		if e := recover(); e != nil {
			p.invokeAdvices(AfterPanic, bean, methodName, args)
		}

		//@After
		err = p.invokeAdvices(After, bean, methodName, args)

		if err == nil && retFunc != nil {
			retFunc()
		}
	}()

	if err != nil {
		//@AfterError
		p.invokeAdvices(AfterError, bean, methodName, args)
	} else {
		//@AfterReturning
		p.invokeAdvices(AfterReturning, bean, methodName, args)
	}

	return
}
