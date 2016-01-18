package aop

type Proxy struct {
	id    string
	funcs map[string]interface{}
}

func NewProxy(beanID string) *Proxy {
	return &Proxy{
		id:    beanID,
		funcs: make(map[string]interface{}),
	}
}

func (p *Proxy) BeanID() string {
	return p.id
}

func (p *Proxy) Method(name string) interface{} {
	fn, _ := p.funcs[name]
	return fn
}

func (p *Proxy) registryFunc(name string, fn interface{}) {
	p.funcs[name] = fn
}
