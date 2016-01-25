package main

import (
	"fmt"

	"github.com/gogap/aop"
)

type TestBean struct {
}

func (p *TestBean) Hello(name string) string {
	fmt.Println("hello", name)
	return "ok"
}

func (p *TestBean) Before(name string) string {
	fmt.Println("before hello", name)
	return "before:I am ok"
}

func (p *TestBean) After() string {
	fmt.Println("after hello")
	return "after:I am ok"
}

type TestBean2 struct {
}

func (p *TestBean2) Foo() {
	fmt.Println("Bar")
}

func main() {
	beanFactory := aop.NewClassicBeanFactory()
	beanFactory.RegisterBean("test_bean", new(TestBean))
	beanFactory.RegisterBean("test_bean2", new(TestBean2))

	gogapAop := aop.NewAOP()

	gogapAop.SetBeanFactory(beanFactory)

	aspect1 := aop.NewAspect("hello", "test_bean")
	aspect1.SetBeanFactory(beanFactory)

	pointcut := aop.NewPointcut("pointcut_1", ".*?")

	aspect1.AddPointcut(pointcut)

	// Before()-> Hello() -> After()
	aspect1.AddAdvice(&aop.Advice{Ordering: aop.Before, Method: "Before", PointcutRefID: "pointcut_1"})
	aspect1.AddAdvice(&aop.Advice{Ordering: aop.After, Method: "After", PointcutRefID: "pointcut_1"})

	gogapAop.AddAspect(aspect1)

	aspect2 := aop.NewAspect("hello2", "test_bean2")
	aspect2.SetBeanFactory(beanFactory)

	// BeforeHello01()-> BeforeHello02()-> Hello() -> AfterHello()
	aspect2.AddAdvice(&aop.Advice{Ordering: aop.After, Method: "Foo", Pointcut: ".*?"})

	gogapAop.AddAspect(aspect2)

	// Get proxy
	proxy, err := gogapAop.GetProxy("test_bean")

	fmt.Println("* Call by Proxy with func type assertion")

	ret := proxy.Method(new(TestBean).Hello).(func(string) string)("I AM Proxy")
	fmt.Println(" -> return value is:", ret)

	fmt.Println("\n* Call by Proxy by Invoke with callback")

	ret2 := ""
	retCallback := func(v string) {
		ret2 = v
	}

	if err = proxy.Invoke(new(TestBean).Hello, "this is params").End(retCallback); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(" -> return value is:", ret2)
	}
}
