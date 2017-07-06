package router

import (
	sysController "GoEjin/system/controller"
	"GoEjin/controller"
)

/**
Route table
 */
var RouteTable = sysController.Router{
	"HOME": sysController.Cfg{&controller.HomeController{}, "", sysController.MethodMap{ "Index" : "" , }},
}
