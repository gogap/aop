AOP
===
Aspect Oriented Programming For Golang

> current version is in alpha, welcome to submit your ideas (api is not stable current version)


### Basic Usage

#### define struct

```go
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

// use join point to around the real func of login
func (p *Auth) Around(pjp aop.ProceedingJoinPointer) {
	fmt.Println("@Begin Around")

	ret := pjp.Proceed("fakeName", "fakePassword")
	ret.MapTo(func(loginResult bool) {
		fmt.Println("@Proceed Result is", loginResult)
	})

	fmt.Println("@End Around")
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
aspect.AddAdvice(&aop.Advice{Ordering: aop.Around, Method: "Around", PointcutRefID: "pointcut_1"})
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
Before Login: zeal
After Login: zeal gogap
Login result: true
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

#### Do not want to assertion func type

```go
proxy.Invoke(new(Auth).Login, "zeal", "errorpassword").End(
		func(result bool) {
			login = result
		})
```

#### Weaving other beans into aspect

##### define a bean

```go
type Foo struct {
}

// @AfterReturning, the method could have args of aop.Result,
// it will get the result from real func return values
func (p *Foo) Bar(result aop.Result) {
	result.MapTo(func(v bool) {
		fmt.Println("Bar Bar Bar .... Result is:", v)
	})
}
```

##### register bean

```go
beanFactory.RegisterBean("foo", new(Foo))
```

##### create aspect
```go
aspectFoo := aop.NewAspect("aspect_2", "foo")
aspectFoo.SetBeanFactory(beanFactory)
```


##### add advice

```go
aspectFoo.AddAdvice(&aop.Advice{Ordering: aop.AfterReturning, Method: "Bar", PointcutRefID: "pointcut_1"})
```

##### add aspect into aop

```go
gogapAop.AddAspect(aspectFoo)
```

result

```bash
Before Login: zeal
Bar Bar Bar .... Result is: true
After Login: zeal gogap
Login result: true
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
go run main.go
==========Func Type Assertion==========
Before Login: zeal
@Begin Around
@Login fakeName fakePassword
@Proceed Result is false
@End Around
After Login: zeal gogap
Login result: false
================Invoke=================
Before Login: zeal
@Begin Around
@Login fakeName fakePassword
@Proceed Result is false
@End Around
After Login: zeal errorpassword
Login result: false
1 aqpk3jjhssa5ul6pt0h0 auth main.(Auth).Login Before
2 aqpk3jjhssa5ul6pt0h0 auth main.(Auth).Login Around
3 aqpk3jjhssa5ul6pt0h0 auth main.(Auth).Login *Login
4 aqpk3jjhssa5ul6pt0h0 foo main.(Auth).Login Bar
5 aqpk3jjhssa5ul6pt0h0 auth main.(Auth).Login After
6 aqpk3jjhssa5ul6pt0hg auth main.(Auth).Login Before
7 aqpk3jjhssa5ul6pt0hg auth main.(Auth).Login Around
8 aqpk3jjhssa5ul6pt0hg auth main.(Auth).Login *Login
9 aqpk3jjhssa5ul6pt0hg foo main.(Auth).Login Bar
10 aqpk3jjhssa5ul6pt0hg auth main.(Auth).Login After
```
> the `*` means the real func in this call