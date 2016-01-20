AOP
===
Aspect Oriented Programming For Golang

> current version is in alpha, welcome to submit your ideas (api is not stable current version)

### Usage:

```go
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

func (p *TestBean) World() {
	fmt.Println("hello", "world")
	return
}

func (p *TestBean) BeforeHello01(name string) string {
	fmt.Println("before hello 01", name)
	return "before:i am ok"
}

func (p *TestBean) BeforeHello02(name string) (err error) {
	fmt.Println("before hello 02", name)
	return nil
}

func (p *TestBean) BeforeWorld() string {
	fmt.Println("before world")
	return "before:i am ok"
}

func (p *TestBean) AfterHello() (err error) {
	fmt.Println("after hello")
	return nil
}

func main() {
	beanFactory := aop.NewClassicBeanFactory()
	beanFactory.RegisterBean("test_bean", "main.TestBean", new(TestBean))

	gogapAop := aop.NewAOP()

	gogapAop.SetBeanFactory(beanFactory)

	aspect := aop.NewAspect("hello", "test_bean")
	aspect.SetBeanFactory(beanFactory)

	// BeforeHello01()-> BeforeHello02()-> Hello() -> AfterHello()
	aspect.AddAdvice(&aop.Advice{Ordering: aop.Before, Method: "BeforeHello01", Pointcut: "Hello()"})
	aspect.AddAdvice(&aop.Advice{Ordering: aop.Before, Method: "BeforeHello02", Pointcut: "Hello()"})
	aspect.AddAdvice(&aop.Advice{Ordering: aop.After, Method: "AfterHello", Pointcut: "Hello()"})

	// BeforeWorld() -> World()
	aspect.AddAdvice(&aop.Advice{Ordering: aop.Before, Method: "BeforeWorld", Pointcut: "World()"})

	gogapAop.AddAspect(aspect)

	// Get proxy
	proxy, err := gogapAop.GetProxy("test_bean")

	fmt.Println("* Call by Proxy with func type assertion")

	ret := proxy.Method("Hello").(func(string) string)("I AM Proxy")
	fmt.Println(" -> return value is:", ret)

	fmt.Println("\n* Call by Proxy by Invoke with callback")

	ret2 := ""
	retCallback := func(v string) {
		ret2 = v
	}

	if err = proxy.Invoke("Hello", "this is params").End(retCallback); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(" -> return value is:", ret2)
	}
}

```

[![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/gogap/aop/trend.png)](https://bitdeli.com/free "Bitdeli Badge")