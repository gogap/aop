package aop

import (
	"reflect"
)

type JoinPointer interface {
	Args() Args
	Target() *Bean
	CallID() string
}

type ProceedingJoinPointer interface {
	JoinPointer

	Proceed(args ...interface{}) InvokeResult
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
	JoinPoint

	method interface{}
}

func (p *ProceedingJoinPoint) Args() Args {
	return p.args
}

func (p *ProceedingJoinPoint) Target() *Bean {
	return p.target
}

func (p *ProceedingJoinPoint) Proceed(args ...interface{}) (ir InvokeResult) {
	v := reflect.ValueOf(p.method)

	var vArgs []reflect.Value
	for _, arg := range args {
		vArgs = append(vArgs, reflect.ValueOf(arg))
	}

	ir.values = v.Call(vArgs)

	return
}
