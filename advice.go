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
	Ordering   AdviceOrdering
	Method     string
	Expression string
}
