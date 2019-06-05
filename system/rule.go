package system

import (
	"fmt"
	"github.com/ejin66/goejin/util"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var routeMap map[string]Cfg

func LoadRouter(table Router) {
	if table == nil {
		return
	}

	routeMap = table
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
			util.PrintError(err)
			io.WriteString(w, util.Error404())
		}
	}()
	uri := req.RequestURI

	fmt.Print(req.Method, " ", req.Host, " ", uri, " ", req.RemoteAddr, " ")

	if uri == "/" {
		io.WriteString(w, util.Web("index.html"))
		return
	}

	data := strings.Split(uri, "/")
	if v, ok := routeMap[strings.ToUpper(data[1])]; ok {
		parse(&v, data[2:], &w, req)
	} else {
		//try to load static resources in src/web/
		util.PrintError("No router found, try to load static resource")
		io.WriteString(w, util.Web(uri[1:]))
	}
}

func parse(cfg *Cfg, data []string, w *http.ResponseWriter, req *http.Request) {

	//New会创建一个指向值的pointer
	//create new instance and the pointer to it : b
	b := reflect.New(reflect.ValueOf(cfg.Cb).Elem().Type()).Interface().(Base)
	ctx := &Context{W: w, Req: req}
	b.Context(ctx)

	//过滤部分请求
	if ok, errMsg := b.Filter(); !ok {
		util.PrintError("Filter out")
		io.WriteString(*w, errMsg)
		return
	}

	if len(data) == 0 {
		data = append(data, "index")
	}

	//将方法名转换成{首字母大写、其余小写}的形式
	methodName := strings.ToUpper(string(data[0][0])) + strings.ToLower(string(data[0][1:]))

	//某些内置方法，这里做特殊判断
	if methodName == "Filter" || methodName == "Context" {
		util.PrintError("Build-in method")
		io.WriteString(*w, util.Error404())
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
		util.PrintError("Request method mismatching")
		io.WriteString(*w, util.Error404())
		return
	}

	mMethod := reflect.ValueOf(b).MethodByName(methodName)
	if mMethod.IsValid() {
		var params []reflect.Value
		if len(data) > 1 {
			params = make([]reflect.Value, len(data)-1)
			for i := range params {
				switch mMethod.Type().In(i).Name() {
				case "int":
					if v, err := strconv.Atoi(data[i+1]); err == nil {
						params[i] = reflect.ValueOf(v)
					} else {
						util.PrintError("Arguments Type is not match!")
						io.WriteString(*w, util.Error404())
						return
					}
				default:
					params[i] = reflect.ValueOf(data[i+1])
				}

			}
		}

		mMethod.Call(params)
		fmt.Println()
	} else {
		util.PrintError("Not found method in controller")
		io.WriteString(*w, util.Error404())
	}

}
