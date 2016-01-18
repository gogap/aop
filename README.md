AOP
===
Aspect Oriented Programming For Golang

> current version is in alpha, welcome to submit your ideas

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

	fmt.Println("Pointcut: Hello()")
	if err := gogapAop.Invoke(
		"test_bean",       // bean id
		"Hello",           // call func
		aop.Args{"gogap"}, // args
		func(ret string) { // the func return value
			fmt.Println("return value is:", ret)
		}); err != nil {
		fmt.Println("call error:", err)
	}

	fmt.Println("")
	fmt.Println("Pointcut: World()")
	if err := gogapAop.Invoke(
		"test_bean", // bean id
		"World",     // call func
		nil,         // args
	); err != nil {
		fmt.Println("call error:", err)
	}
}

```

```bash
$> go run main.go

Pointcut: Hello()
before hello 01 gogap
before hello 02 gogap
hello gogap
after hello
return value is: ok

Pointcut: World()
before world
hello world
```