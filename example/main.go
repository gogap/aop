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

func main() {
	beanFactory := aop.NewClassicBeanFactory()
	beanFactory.RegisterBean("test_bean", new(TestBean))

	gogapAspect := aop.NewAspect("gogap_aop")

	gogapAspect.SetBeanFactory(beanFactory)

	gogapAspect.Invoke(
		"test_bean",       // bean id
		"Hello",           // call func
		aop.Args{"gogap"}, // args
		func(ret string) { // the func return value
			fmt.Println("return value is:", ret)
		})
}
