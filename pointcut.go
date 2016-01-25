package aop

import (
	"regexp"
)

type Pointcut struct {
	ID         string
	Expression string
}

func NewPointcut(id, expr string) *Pointcut {
	return &Pointcut{ID: id, Expression: expr}
}

func (p *Pointcut) IsMatch(bean *Bean, methodName string, args Args) (matched bool, err error) {
	return regexp.MatchString(p.Expression, methodName)
}
