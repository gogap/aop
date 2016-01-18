package aop

import (
	"github.com/gogap/errors"
)

const (
	AOPErrorNamespace = "AOP"
)

var (
	ErrBeanInstanceIsNil        = errors.TN(AOPErrorNamespace, 1, "aop error namespace is nil, id: {{.id}}")
	ErrBeanIsNotAnPtr           = errors.TN(AOPErrorNamespace, 2, "bean should be an ptr, id: {{.id}}")
	ErrBeanAlreadyRegistered    = errors.TN(AOPErrorNamespace, 3, "bean already regitered, id: {{.id}}")
	ErrBeanIDShouldNotBeEmpty   = errors.TN(AOPErrorNamespace, 4, "bean id should not be empty")
	ErrBeanNotExist             = errors.TN(AOPErrorNamespace, 5, "bean not exist, id: {{.id}}")
	ErrAspectIDShouldNotBeEmpty = errors.TN(AOPErrorNamespace, 6, "aspect id should not be empty")
	ErrBeanMethodNotExit        = errors.TN(AOPErrorNamespace, 7, "bean method not exist, id: {{.id}}, class: {{.class}}, method: {{.method}}")
	ErrWrongAdviceFuncArgsNum   = errors.TN(AOPErrorNamespace, 8, "wrong advice func args number, id: {{.id}}, class: {{.class}}, method: {{.method}}")
)
