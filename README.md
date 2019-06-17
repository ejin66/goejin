# Introduce

**It route the url like :**

```go 
/controller/function/arg1/arg2... 
```

**to controller and call the function of the same name , like :**

```go 
Controller.function(arg1,arg2,...)
```


# How To Use

### 1. 下载
```bash
go get github.com/ejin66/goejin
```

### 2. 配置参数
创建config.json，设置ip,port,db等参数，参考[example/config.json](https://github.com/ejin66/goejin/blob/master/example/config.json).

### 3. 新建controller
创建一个内嵌有system.BaseController的struct, 参照[example/home.go](https://github.com/ejin66/goejin/blob/master/example/home.go)中的HomeController:

```go
type HomeController struct {
	controller.BaseController
}
```

### 4. 配置Router
配置路由规则，如：

```go
var RouteTable = system.Router{
	"home": system.Cfg{&HomeController{}, "", system.MethodMap{ "Index":"GET"}},
}
//  .../home/index  --> HomeController的index方法
//  且必须 request method 限制为GET
```

其中， key值是url中的controller名， value值：第一个参是新建的controller实例指针； 第二个参是默认方法 POST/GET/"" ,代表默认请求方法,空字符串是无限制；第三个参 是限制特定方法名的请求方法。

### 5. 运行
```go
//第一个参数： 第二步中创建的config.json路径
//第二个参数： 第四步中配置的路由
goejin.Listen("./config.json", RouteTable)
```

# Example
框架中有example代码，可以直接跑起来看效果。
起服务：
```bash
cd example
go run ./
```
接着，访问：
- `http://localhost/home/index`
- `http://localhost/`

可查看效果。
