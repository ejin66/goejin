package router

import (
	"net/http"
	"io"
	"fmt"
	"strings"
	"reflect"
	"api/controller"
	"api/config"
)

var routeMap map[string]config.Cfg

func init() {
	routeMap = config.ROUTER
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
			fmt.Println("get error:", err)
			io.WriteString(w, config.Error404())
		}
	}()
	uri := req.RequestURI

	fmt.Println(req.Method , "request ", uri)

	var b controller.Base
	if uri == "/" {
		uri = config.HOME_URI
	}
	data := strings.Split(uri, "/")
	if v, ok := routeMap[strings.ToUpper(data[1])]; ok {
		b = v.Cb.Instance()
		fmt.Printf("%p\n", b)
		ctx := &controller.Context{&w, req}
		b.Context(ctx)
		fmt.Printf("%p\n", b)

		//若请求方法与配置的方法不一致，返回404
		if v.Method != "" && req.Method != v.Method {
			io.WriteString(w, config.Error404())
			return
		}
	}
	parse(b, data[2:], &w, req)
}

func parse(b controller.Base, data []string, w *http.ResponseWriter, req *http.Request) {
	if b != nil {
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

		mMethod := reflect.ValueOf(b).MethodByName(methodName)

		if mMethod.IsValid() {
			mMethod.Call(params)
		} else {
			io.WriteString(*w, config.Error404())
		}

	} else {
		io.WriteString(*w, config.Error404())
	}
}


