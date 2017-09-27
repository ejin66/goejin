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

**1. 配置参数。在config/config.json中配置，如ip,port,db相关参数等等**

**2. 新建controller。 在controller/下创建一个内嵌有controller.BaseController的struct, 参照controller/home.go中的HomeController**

**3. 配置Router。在router/router.go中，配置路由规则，如：**

```go
var RouteTable = sysController.Router{
	"HOME": sysController.Cfg{&controller.HomeController{}, "", sysController.MethodMap{ "Index":"GET"}}
}
//  .../HOME/index  --> HomeController的index方法
//  且必须 request method 限制为GET
```

其中， key值是url中的controller名， value值：第一个参是新建的controller实例指针； 第二个参是默认方法 POST/GET/"" ,代表默认请求方法,空字符串是无限制；第三个参 是限制特定方法的请求。


以上
