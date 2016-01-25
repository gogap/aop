package aop

import (
	"regexp"
	"strings"

	"github.com/gogap/errors"
)

type Pointcut struct {
	id         string
	expression string

	execExpr   string
	withinExpr string
}

func NewPointcut(id, expr string) *Pointcut {
	regex := regexp.MustCompile(`execution\s?\((.*)\(\)\)`)

	execExprs := regex.FindStringSubmatch(expr)
	if len(execExprs) != 2 {
		err := ErrBadPointcutExpr.New(errors.Params{"expr": expr})
		panic(err)
	}
	return &Pointcut{id: id, expression: expr, execExpr: execExprs[1]}
}

func (p *Pointcut) ID() string {
	return p.id
}

func (p *Pointcut) Expression() string {
	return p.expression
}

func (p *Pointcut) IsMatch(bean *Bean, methodName string, args Args) (matched bool, err error) {
	// match execution
	return p.isExecExprMatch(bean, methodName, args)
}

func (p *Pointcut) isExecExprMatch(bean *Bean, methodName string, args Args) (matched bool, err error) {
	if methodName == p.execExpr {
		return true, nil
	} else {
		indexExprPkg := strings.LastIndex(p.execExpr, ".")
		indexMethodPkg := strings.LastIndex(methodName, ".")

		expr := p.execExpr
		if indexExprPkg < 0 {
			expr = methodName[:indexMethodPkg] + "." + expr
			indexExprPkg = indexMethodPkg
		}

		execExprPkg := expr[:indexExprPkg]
		methodPkg := methodName[:indexMethodPkg]

		pkgExprMatch := false
		if pkgExprMatch, err = regexp.MatchString(execExprPkg, methodPkg); err != nil {
			return
		} else if pkgExprMatch {
			return true, nil
		}

		if execExprPkg == methodPkg ||
			pkgExprMatch == true {

			execExprFn := expr[strings.LastIndex(expr, ".")+1:]
			methodFn := methodName[strings.LastIndex(methodName, ".")+1:]

			fnExprMatch := false
			if fnExprMatch, err = regexp.MatchString(execExprFn, methodFn); err != nil {
				return
			}

			if fnExprMatch || execExprFn == methodFn {
				return true, nil
			}
		}

		return false, nil
	}
	return false, nil
}
