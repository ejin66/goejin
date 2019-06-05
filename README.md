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

**1. 下载**
```bash
go get github.com/ejin66/goejin
# goejin中封装了db操作，需要依赖mysql driver
go get github.com/go-sql-driver/mysql
```

**2. 配置参数。**
创建config.json，设置ip,port,db等参数，参考example/config.json.

**3. 新建controller。**
创建一个内嵌有system.BaseController的struct, 参照example/homeController.go中的HomeController:

```go
type HomeController struct {
	controller.BaseController
}
```

**3. 配置Router。**
在router/router.go中，配置路由规则，如：

```go
var RouteTable = sysController.Router{
	"home": sysController.Cfg{&controller.HomeController{}, "", sysController.MethodMap{ "Index":"GET"}}
}
//  .../home/index  --> HomeController的index方法
//  且必须 request method 限制为GET
```

其中， key值是url中的controller名， value值：第一个参是新建的controller实例指针； 第二个参是默认方法 POST/GET/"" ,代表默认请求方法,空字符串是无限制；第三个参 是限制特定方法的请求。


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
