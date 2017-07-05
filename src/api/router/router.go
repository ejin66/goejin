package router

import (
	"api/controller"
	"api/system"
)

/**
Route table
 */
var RouteTable = system.Router{
	"HOME": system.Cfg{&controller.HomeController{}, "", system.MethodMap{ "Index" : "" , }},
}
