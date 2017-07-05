package router

import (
	"net/http"
	"io"
	"fmt"
	"strings"
	"reflect"
	"api/config"
	"api/system"
)

var routeMap map[string]system.Cfg

func init() {
	routeMap = config.RouteTable
}

//路由规则
func GetServeMux() *http.ServeMux {
	var mServeMux http.ServeMux
	mServeMux.HandleFunc("/", defaultHandler)
	return &mServeMux
}

//路由规则
func defaultHandler(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			//这里，主要是捕获调用函数参数不一致情况
			system.PrintError(err)
			io.WriteString(w, system.Error404())
		}
	}()
	uri := req.RequestURI

	fmt.Println(req.Method , "request ", uri)

	if uri == "/" {
		uri = config.HOME_URI
	}
	data := strings.Split(uri, "/")
	if v, ok := routeMap[strings.ToUpper(data[1])]; ok {
		parse(&v, data[2:], &w, req)
	} else {
		io.WriteString(w, system.Error404())
	}
}

func parse(cfg *system.Cfg, data []string, w *http.ResponseWriter, req *http.Request) {
	b := cfg.Cb.Instance()
	//fmt.Printf("%p\n", b)
	ctx := &system.Context{w, req}
	b.Context(ctx)
	//fmt.Printf("%p\n", b)

	if len(data) == 0 {
		data = append(data, "index")
	}
	var params []reflect.Value
	if len(data) > 1 {
		params = make([]reflect.Value, len(data)-1)
		for i, _ := range params {
			params[i] = reflect.ValueOf(data[i+1])
		}
	}
	//将方法名转换成{首字母大写、其余小写}的形式
	methodName := strings.ToUpper(string(data[0][0])) + strings.ToLower(string(data[0][1:]))

	if methodName == "Instance" {
		io.WriteString(*w, system.Error404())
		return
	}

	//根据配置文件，得到配置的请求方法
	requestMethod := cfg.DefaultMethod
	for key, value := range cfg.Methods {
		if key == methodName {
			requestMethod = value
			break
		}
	}

	//若请求方法不一致，直接抛出error
	if requestMethod != "" && requestMethod != req.Method {
		io.WriteString(*w, system.Error404())
		return
	}

	mMethod := reflect.ValueOf(b).MethodByName(methodName)

	if mMethod.IsValid() {
		mMethod.Call(params)
	} else {
		io.WriteString(*w, system.Error404())
	}

}


