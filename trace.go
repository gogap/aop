package aop

import (
	"sync"
	"sync/atomic"
	"time"
)

var (
	trace       *Trace
	isTracing   bool
	traceLocker sync.Mutex
)

type TraceItem struct {
	ID        int
	InvokeID  string
	File      string
	Line      int
	Method    string
	Pointcut  string
	BeanRefID string
	Timestamp string
}

type Trace struct {
	id    int64
	items []TraceItem
}

func newTrace() *Trace {

	trace = new(Trace)
	trace.id = 0

	return trace
}

func (p *Trace) append(invokeID, file string, line int, method, pointcut, beanRefID string) {
	item := TraceItem{
		ID:        int(atomic.AddInt64(&p.id, 1)),
		File:      file,
		Line:      line,
		Method:    method,
		Pointcut:  pointcut,
		InvokeID:  invokeID,
		BeanRefID: beanRefID,
		Timestamp: time.Now().Format(time.StampNano),
	}

	p.items = append(p.items, item)
}

func (p *Trace) Items() []TraceItem {
	return p.items
}

func StartTrace() (err error) {
	traceLocker.Lock()
	defer traceLocker.Unlock()

	if isTracing {
		err = ErrTracAlreadyStarted.New()
		return
	}

	isTracing = true
	trace = newTrace()

	return
}

func StopTrace() (t *Trace, err error) {
	traceLocker.Lock()
	defer traceLocker.Unlock()

	if !isTracing {
		err = ErrTracNotStarted.New()
		return
	}

	isTracing = false
	t = trace

	return
}

func IsTracing() bool {
	return isTracing
}

func appendTraceItem(invokeID, file string, line int, method, pointcut, beanRefID string) {
	traceLocker.Lock()
	defer traceLocker.Unlock()

	if !isTracing {
		return
	}

	trace.append(invokeID, file, line, method, pointcut, beanRefID)
}
