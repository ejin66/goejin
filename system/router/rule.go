package router

import (
	"net/http"
	"io"
	"fmt"
	"strings"
	"reflect"
	"GoEjin/router"
	"GoEjin/system/common"
	"GoEjin/system/controller"
	"GoEjin/system/config"
)

var routeMap map[string]controller.Cfg

func init() {
	routeMap = router.RouteTable
}

func GetServeMux() *http.ServeMux {
	var mServeMux http.ServeMux
	mServeMux.HandleFunc("/", defaultHandler)
	return &mServeMux
}

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			//这里，主要是捕获调用函数参数不一致情况
			common.PrintError(err)
			io.WriteString(w, common.Error404())
		}
	}()
	uri := req.RequestURI


	fmt.Print(req.Method, " ", req.Host, " ", uri, " ",req.RemoteAddr, " " )

	if uri == "/" {
		uri = "/" + config.GetConfig().DEFAULT_CONTROLLER
	}
	data := strings.Split(uri, "/")
	if v, ok := routeMap[strings.ToUpper(data[1])]; ok {
		parse(&v, data[2:], &w, req)
	} else {
		common.PrintError("No router found")
		io.WriteString(w, common.Error404())
	}
}

func parse(cfg *controller.Cfg, data []string, w *http.ResponseWriter, req *http.Request) {

	//create new instance and the pointer to it : b
	b := reflect.New(reflect.ValueOf(cfg.Cb).Elem().Type()).Interface().(controller.Base)
	ctx := &controller.Context{w, req}
	b.Context(ctx)

	//过滤部分请求
	if ok, errMsg := b.Filter(); !ok {
		common.PrintError("Filter out")
		io.WriteString(*w, errMsg)
		return
	}

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

	//某些内置方法，这里做特殊判断
	if methodName == "Filter" || methodName == "Context" {
		common.PrintError("Build-in method")
		io.WriteString(*w, common.Error404())
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
		common.PrintError("Request method mismatching")
		io.WriteString(*w, common.Error404())
		return
	}

	mMethod := reflect.ValueOf(b).MethodByName(methodName)

	if mMethod.IsValid() {
		mMethod.Call(params)
		fmt.Println()
	} else {
		common.PrintError("Not found method in controller")
		io.WriteString(*w, common.Error404())
	}

}
