package aop

import (
	"regexp"
	"runtime"
	"strings"
	"sync"
)

type Pointcut struct {
	id string

	executionExprs    []*regexp.Regexp
	notExecutionExprs []*regexp.Regexp
	byBeans           []*regexp.Regexp
	notByBeans        []*regexp.Regexp
	withIn            []*regexp.Regexp
	notWithIn         []*regexp.Regexp

	exprLocker sync.Mutex
}

func (p *Pointcut) Execution(expr string) *Pointcut {
	p.exprLocker.Lock()
	defer p.exprLocker.Unlock()

	p.executionExprs = append(p.executionExprs, regexp.MustCompile(expr))
	return p
}

func (p *Pointcut) NotExecution(expr string) *Pointcut {
	p.exprLocker.Lock()
	defer p.exprLocker.Unlock()

	p.notExecutionExprs = append(p.notExecutionExprs, regexp.MustCompile(expr))
	return p
}

func (p *Pointcut) Bean(expr string) *Pointcut {
	p.exprLocker.Lock()
	defer p.exprLocker.Unlock()

	p.byBeans = append(p.byBeans, regexp.MustCompile(expr))
	return p
}

func (p *Pointcut) NotBean(expr string) *Pointcut {
	p.exprLocker.Lock()
	defer p.exprLocker.Unlock()

	p.notByBeans = append(p.notByBeans, regexp.MustCompile(expr))
	return p
}

func (p *Pointcut) Within(expr string) *Pointcut {
	p.exprLocker.Lock()
	defer p.exprLocker.Unlock()

	p.withIn = append(p.withIn, regexp.MustCompile(expr))
	return p
}

func (p *Pointcut) NotWithin(expr string) *Pointcut {
	p.exprLocker.Lock()
	defer p.exprLocker.Unlock()

	p.notWithIn = append(p.notWithIn, regexp.MustCompile(expr))
	return p
}

func NewPointcut(id string) *Pointcut {
	return &Pointcut{
		id: id,
	}
}

func (p *Pointcut) ID() string {
	return p.id
}

func (p *Pointcut) IsMatch(bean *Bean, methodName string, args Args) bool {

	// not execution
	for _, notExecRegex := range p.notExecutionExprs {
		if notExecRegex.MatchString(methodName) {
			return false
		}
	}

	// execution
	execGot := false
	for _, execRegex := range p.executionExprs {
		if execRegex.MatchString(methodName) {
			execGot = true
			break
		}
	}

	if !execGot {
		return false
	}

	// not bean
	for _, notBeanRegex := range p.notByBeans {
		if notBeanRegex.MatchString(bean.id) {
			return false
		}
	}

	// bean
	if len(p.byBeans) > 0 {
		beanGot := false
		for _, beanRegex := range p.byBeans {
			if beanRegex.MatchString(bean.id) {
				beanGot = true
				break
			}
		}

		if !beanGot {
			return false
		}
	}

	stacks := ""
	if len(p.notWithIn) > 0 || len(p.withIn) > 0 {
		buf := make([]byte, 4096)
		runtime.Stack(buf, false)
		stacks = strings.Replace(string(buf), "*", "", -1)
	}

	// not within
	for _, notWithInRegex := range p.notWithIn {
		if notWithInRegex.MatchString(stacks) {
			return false
		}
	}

	// with in
	if len(p.withIn) > 0 {
		withInGot := false
		for _, withInRegex := range p.withIn {
			if withInRegex.MatchString(stacks) {
				withInGot = true
				break
			}
		}

		if !withInGot {
			return false
		}
	}

	return true
}
