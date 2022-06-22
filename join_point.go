package aop

import (
	"reflect"
)

type joinPointFunc func(JoinPointer) error
type joinPointWithResultFunc func(JoinPointer, Result) error
type proceedingJoinPoint func(ProceedingJoinPointer) error

var (
	joinPointFuncType           = reflect.TypeOf((*joinPointFunc)(nil)).Elem()
	joinPointWithResultFuncType = reflect.TypeOf((*joinPointWithResultFunc)(nil)).Elem()
	proceedingJoinPointType     = reflect.TypeOf((*proceedingJoinPoint)(nil)).Elem()
)

var (
	_ JoinPointer           = (*JoinPoint)(nil)
	_ ProceedingJoinPointer = (*ProceedingJoinPoint)(nil)
)

type JoinPointer interface {
	Args() Args
	Target() *Bean
	CallID() string
}

type ProceedingJoinPointer interface {
	JoinPointer
	Proceed(args ...interface{}) Result
}

type JoinPoint struct {
	callID string
	args   Args
	target *Bean
}

func (p *JoinPoint) CallID() string {
	return p.callID
}

func (p *JoinPoint) Args() Args {
	return p.args
}

func (p *JoinPoint) Target() *Bean {
	return p.target
}

type ProceedingJoinPoint struct {
	result []reflect.Value
	JoinPointer
	method func(args ...interface{}) []reflect.Value
}

func (p *ProceedingJoinPoint) Args() Args {
	return p.JoinPointer.Args()
}

func (p *ProceedingJoinPoint) Target() *Bean {
	return p.JoinPointer.Target()
}

func (p *ProceedingJoinPoint) Proceed(args ...interface{}) Result {
	p.result = p.method(args...)
	return p.result
}
