package aop

import (
	"github.com/gogap/errors"
)

const (
	AOPErrorNamespace = "AOP"
)

var (
	ErrBeanInstanceIsNil         = errors.TN(AOPErrorNamespace, 1, "aop error namespace is nil, id: {{.id}}")
	ErrBeanIsNotAnPtr            = errors.TN(AOPErrorNamespace, 2, "bean should be an ptr, id: {{.id}}")
	ErrBeanAlreadyRegistered     = errors.TN(AOPErrorNamespace, 3, "bean already regitered, id: {{.id}}")
	ErrBeanIDShouldNotBeEmpty    = errors.TN(AOPErrorNamespace, 4, "bean id should not be empty")
	ErrBeanNotExist              = errors.TN(AOPErrorNamespace, 5, "bean not exist, id: {{.id}}")
	ErrAspectIDShouldNotBeEmpty  = errors.TN(AOPErrorNamespace, 6, "aspect id should not be empty")
	ErrBeanMethodNotExit         = errors.TN(AOPErrorNamespace, 7, "bean method not exist, id: {{.id}}, class: {{.class}}, method: {{.method}}")
	ErrWrongAdviceFuncArgsNum    = errors.TN(AOPErrorNamespace, 8, "wrong advice func args number, id: {{.id}}, class: {{.class}}, method: {{.method}}")
	ErrEndInvokeParamsIsNotFunc  = errors.TN(AOPErrorNamespace, 9, "en invoke params is not func, bean id: {{.id}}, method: {{.method}}")
	ErrWrongEndInvokeFuncArgsNum = errors.TN(AOPErrorNamespace, 10, "wrong end invoke func args number, bean id: {{.id}}, method: {{.method}}")
	ErrInvokeParamsIsNotFunc     = errors.TN(AOPErrorNamespace, 11, "invoke params is not func, bean id: {{.id}}, method: {{.method}}")
	ErrWrongInvokeFuncArgsNum    = errors.TN(AOPErrorNamespace, 12, "wrong invoke func args number, bean id: {{.id}}, method: {{.method}}")
	ErrInvokeFuncNotExist        = errors.TN(AOPErrorNamespace, 13, "invoke func not exist, bean id: {{.id}}, method: {{.method}}")
	ErrInvokeFuncTypeError       = errors.TN(AOPErrorNamespace, 14, "invoke func is not func type, bean id: {{.id}}, method: {{.method}}")
	ErrEndInvokeTwice            = errors.TN(AOPErrorNamespace, 15, "end invoke twice, bean id: {{.id}}, method: {{.method}}")
	ErrBadInvokeMethodType       = errors.TN(AOPErrorNamespace, 16, "invoke method params should be func name or func type")
	ErrPointcutNotExist          = errors.TN(AOPErrorNamespace, 17, "pointcut no exist, pointcut id: {{.id}}")
	ErrMethodIsNil               = errors.TN(AOPErrorNamespace, 18, "method is nil")
	ErrBadMethodType             = errors.TN(AOPErrorNamespace, 19, "method type error")
	ErrBadPointcutExpr           = errors.TN(AOPErrorNamespace, 20, "bad pointcut expression, expression: {{.expr}}")
	ErrEmptyPointcutRefID        = errors.TN(AOPErrorNamespace, 21, "pointcut ref id is empty")
	ErrExecExpressionAlreadySet  = errors.TN(AOPErrorNamespace, 22, "execution expression already set at pointcut {{.id}}, old expr is: {{.expr}}")
	ErrTracAlreadyStarted        = errors.TN(AOPErrorNamespace, 23, "trace already started")
	ErrTracNotStarted            = errors.TN(AOPErrorNamespace, 24, "trace not started")
	ErrUnknownJoinPointArgType   = errors.TN(AOPErrorNamespace, 25, "unknown advice func arg type, bean id: {{.id}}, method: {{.method}}, pointcut {{.refID}}")
	ErrJoinPointArgsUsage        = errors.TN(AOPErrorNamespace, 26, "only AfterReturning could have aop.Result argument, bean id: {{.id}}, method: {{.method}}, pointcut {{.refID}}")
	ErrMapperArgShouldBeFunc     = errors.TN(AOPErrorNamespace, 27, "result map to arg should be a func")
	ErrWrongMapFuncArgsNum       = errors.TN(AOPErrorNamespace, 28, "wrong end invoke func args number")
)
