package main

import (
	"fmt"

	"github.com/gogap/aop"
)

type Auth struct {
}

func (p *Auth) Login(userName, password string) bool {
	if userName == "zeal" && password == "gogap" {
		return true
	}
	return false
}

// use join point to get Args from real method
func (p *Auth) Before(jp aop.JoinPointer) {
	username := ""
	jp.Args().MapTo(func(u, p string) {
		username = u
	})

	fmt.Printf("Before Login: %s\n", username)
}

// the args is same as Login
func (p *Auth) After(username, password string) {
	fmt.Printf("After Login: %s %s\n", username, password)
}

type Foo struct {
}

// @AfterReturning, the method could have args of aop.Result,
// it will get the result from real func return values
func (p *Foo) Bar(result aop.Result) {
	result.MapTo(func(v bool) {
		fmt.Println("Bar Bar Bar .... Result is:", v)
	})
}

func main() {
	beanFactory := aop.NewClassicBeanFactory()
	beanFactory.RegisterBean("auth", new(Auth))
	beanFactory.RegisterBean("foo", new(Foo))

	aspect := aop.NewAspect("aspect_1", "auth")
	aspect.SetBeanFactory(beanFactory)

	aspectFoo := aop.NewAspect("aspect_2", "foo")
	aspectFoo.SetBeanFactory(beanFactory)

	pointcut := aop.NewPointcut("pointcut_1").Execution(`Login()`)
	pointcut.Execution(`Login()`)

	aspect.AddPointcut(pointcut)
	aspectFoo.AddPointcut(pointcut)

	aspect.AddAdvice(&aop.Advice{Ordering: aop.Before, Method: "Before", PointcutRefID: "pointcut_1"})
	aspect.AddAdvice(&aop.Advice{Ordering: aop.After, Method: "After", PointcutRefID: "pointcut_1"})
	aspectFoo.AddAdvice(&aop.Advice{Ordering: aop.AfterReturning, Method: "Bar", PointcutRefID: "pointcut_1"})

	gogapAop := aop.NewAOP()
	gogapAop.SetBeanFactory(beanFactory)
	gogapAop.AddAspect(aspect)
	gogapAop.AddAspect(aspectFoo)

	var err error
	var proxy *aop.Proxy

	// Get proxy
	if proxy, err = gogapAop.GetProxy("auth"); err != nil {
		fmt.Println("get proxy failed", err)
		return
	}

	// start trace for debug
	aop.StartTrace()

	fmt.Println("==========Func Type Assertion==========")

	login := proxy.Method(new(Auth).Login).(func(string, string) bool)("zeal", "gogap")

	fmt.Println("Login result:", login)

	fmt.Println("================Invoke=================")

	if err = proxy.Invoke(new(Auth).Login, "zeal", "errorpassword").End(
		func(result bool) {
			login = result
		}); err != nil {
		fmt.Println("invoke proxy func error", err)
	} else {
		fmt.Println("Login result:", login)
	}

	t, _ := aop.StopTrace()

	// print trace result
	for _, item := range t.Items() {
		fmt.Println(item.ID, item.InvokeID, item.BeanRefID, item.Pointcut, item.Method)
	}
}
