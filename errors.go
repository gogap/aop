package aop

import (
	"github.com/gogap/errors"
)

const (
	AOPErrorNamespace = "AOP"
)

var (
	ErrBeanInstanceIsNil      = errors.TN(AOPErrorNamespace, 1, "aop error namespace is nil, id: {{.id}}")
	ErrBeanIsNotAnPtr         = errors.TN(AOPErrorNamespace, 2, "bean should be an ptr, id: {{.id}}")
	ErrBeanAlreadyRegistered  = errors.TN(AOPErrorNamespace, 3, "bean already regitered, id: {{.id}}")
	ErrBeanIDShouldNotBeEmpty = errors.TN(AOPErrorNamespace, 4, "bean id should not be empty")
	ErrBeanNotExist           = errors.TN(AOPErrorNamespace, 5, "bean not exist, id: {{.id}}")
)
