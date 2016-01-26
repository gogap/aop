package main

import (
	"fmt"

	"github.com/gogap/aop"
)

type Auth struct {
}

func (p *Auth) Login(userName string, password string) bool {
	if userName == "zeal" && password == "gogap" {
		return true
	}
	return false
}

func (p *Auth) Before(username string, password string) {
	fmt.Println(username, "begin login")
}

func (p *Auth) After(username string, password string) {
	fmt.Println(username, "logged in")
}

func main() {
	beanFactory := aop.NewClassicBeanFactory()
	beanFactory.RegisterBean("auth", new(Auth))

	aspect := aop.NewAspect("aspect_1", "auth")
	aspect.SetBeanFactory(beanFactory)

	pointcut := aop.NewPointcut("pointcut_1").Execution(`Login()`)

	aspect.AddPointcut(pointcut)

	aspect.AddAdvice(&aop.Advice{Ordering: aop.Before, Method: "Before", PointcutRefID: "pointcut_1"})
	aspect.AddAdvice(&aop.Advice{Ordering: aop.After, Method: "After", PointcutRefID: "pointcut_1"})

	gogapAop := aop.NewAOP()
	gogapAop.SetBeanFactory(beanFactory)
	gogapAop.AddAspect(aspect)

	// Get proxy
	proxy, _ := gogapAop.GetProxy("auth")

	// start trace for debug
	aop.StartTrace()

	login := proxy.Method(new(Auth).Login).(func(string, string) bool)("zeal", "gogap")

	fmt.Println("login result:", login)

	t, _ := aop.StopTrace()

	// print trace result
	for _, item := range t.Items() {
		fmt.Println(item.ID, item.InvokeID, item.BeanRefID, item.Pointcut, item.Method)
	}
}
