package router

import (
	sysController "api/system/controller"
	"api/controller"
)

/**
Route table
 */
var RouteTable = sysController.Router{
	"HOME": sysController.Cfg{&controller.HomeController{}, "", sysController.MethodMap{ "Index" : "" , }},
}
