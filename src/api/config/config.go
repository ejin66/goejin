package config

import (
	"api/controller"
	"api/system"
)

/**
配置ip , port ,home address
 */
const (
	IP_ADDRESS = "127.0.0.1"
	IP_PORT    = "8080"
	HOME_URI   = "/home"
)

/**
配置路由关系表
 */
var RouteTable = system.Router{
	"HOME": system.Cfg{&controller.HomeController{}, "GET", system.MethodMap{ "Index" : "GET" , }},
}
