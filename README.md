AOP
===
Aspect Oriented Programming For Golang

> current version is in alpha, welcome to submit your ideas (api is not stable current version)


### Basic Usage

#### define struct

```go
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
```

In this case, we want call `Before()` func before `Login()`, and `After()` func after `Login()`

In general, we will do it like as following

```go
func (p *Auth) Login(userName string, password string) bool {
	p.Before(userName, password)
	defer p.After(userName, password)
	
	if userName == "zeal" && password == "gogap" {
		return true
	}
	return false
}
```

So, if we have more funcs to call before and after, it will pollution the real logic func `Login()`, we want a proxy help us to invoke `Before()` and `After()` automatic.

That was what AOP does.
 
#### Step 1: Define Beans factory

```go
beanFactory := aop.NewClassicBeanFactory()
beanFactory.RegisterBean("auth", new(Auth))
```
 
#### Step 2: Define Aspect

```go
aspect := aop.NewAspect("aspect_1", "auth")
aspect.SetBeanFactory(beanFactory)
``` 

#### Step 3: Define Pointcut

```go
pointcut := aop.NewPointcut("pointcut_1").Execution(`Login()`)
aspect.AddPointcut(pointcut)
``` 

#### Step 4: Add Advice

```go
aspect.AddAdvice(&aop.Advice{Ordering: aop.Before, Method: "Before", PointcutRefID: "pointcut_1"})
aspect.AddAdvice(&aop.Advice{Ordering: aop.After, Method: "After", PointcutRefID: "pointcut_1"})
```

#### Step 5: Create AOP

```go
gogapAop := aop.NewAOP()
gogapAop.SetBeanFactory(beanFactory)
gogapAop.AddAspect(aspect)
```

#### Setp 6: Get Proxy

```go
proxy, err := gogapAop.GetProxy("auth")
```

#### Last Step: Enjoy

```go
login := proxy.Method(new(Auth).Login).(func(string, string) bool)("zeal", "gogap")

fmt.Println("login result:", login)
```
> output

```bash
$> go run main.go
zeal begin login
zeal logged in
login result: true
```

### Advance

#### Pointcut expression

> every condition expression is regex expression

```go
pointcut := aop.NewPointcut("pointcut_1")

// will trigger the advice while call login
pointcut.Execution(`Login()`)

// will trigger the advice will call any func
pointcut.Execution(`.*?`)

// will not trigger the advice will call any func
pointcut.NotExecution(`Login()`)
```

##### other conditions:
- WithIn
- NotWithIn
- Bean
- NotBean

```go
// will trigger the advie while we call Login 
// and in bean named auth
pointcut.Execution(`Login()`).Bean(`auth`)

// will trigger the advie while we call Login 
// and in bean named auth and sysAuth
pointcut.Execution(`Login()`).Bean(`auth`).Bean(`sysAuth`)


// will trigger the advie while we call Login 
// and in bean named auth not sysAuth
pointcut.Execution(`Login()`).Bean(`auth`).NotBean(`sysAuth`)

// will trigger the advie while we call Login 
// and the call stacktrace should contain example/aop/main
pointcut.Execution(`Login()`).WithIn(`example/aop/main`)

```

#### Turn on trace for debug

```go
err := aop.StartTrace()

....
// use proxy to call your funcs

t, err := aop.StopTrace()

for _, item := range t.Items() {
		fmt.Println(item.ID, item.InvokeID, item.BeanRefID, item.Pointcut, item.Method)
}
```

```bash
$> go run main.go
zeal begin login
zeal logged in
login result: true
1 aqjoq1jhssa4c7sm7a20 auth main.(Auth).Login Before
2 aqjoq1jhssa4c7sm7a20 auth main.(Auth).Login *Login
3 aqjoq1jhssa4c7sm7a20 auth main.(Auth).Login After
```
> the `*` means the real func in this call