package config

import "api/controller"

type Cfg struct {
	Method string
	Cb controller.Base
}

type router map[string]Cfg

/**
配置ip , port ,home address
 */
const (
	IP_ADDRESS  = "127.0.0.1"
	IP_PORT  = "8080"
	HOME_URI = "/home"
)

/**
配置路由关系表
 */
var ROUTER =  router {
	"HOME" : Cfg {"POST",&controller.HomeController{} },
}
