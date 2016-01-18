package aop

type AdviceOrdering int

const (
	Before         AdviceOrdering = 1
	After          AdviceOrdering = 2
	AfterReturning AdviceOrdering = 3
	AfterError     AdviceOrdering = 4
	AfterPanic     AdviceOrdering = 5
	Around         AdviceOrdering = 6
)

type Advice struct {
	Ordering    AdviceOrdering
	Method      string
	Pointcut    string
	PointcutRef *Pointcut

	beanRef *Bean
}

func (p *Advice) IsMatch(ordering AdviceOrdering, bean *Bean, methodName string, args Args) (isMatch bool, err error) {
	if !(ordering == Before && p.Ordering == Around) &&
		!(ordering == After && p.Ordering == Around) &&
		ordering != p.Ordering {
		return
	}

	if p.Pointcut != "" {
		if methodName+"()" == p.Pointcut {
			return true, nil
		} else {
			return false, nil
		}
	}

	if p.PointcutRef != nil {
		if methodName+"()" == p.PointcutRef.Expression {
			return true, nil
		} else {
			return false, nil
		}
	}

	return
}
