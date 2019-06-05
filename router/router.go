package router

import (
	"GoEjin/controller"
	sysController "GoEjin/system/controller"
)

/**
Route table
POST/GET/EMPTY STRING(both post and get)
*/
var RouteTable = sysController.Router{
	"HOME": sysController.Cfg{&controller.HomeController{}, "", sysController.MethodMap{"Index": ""}},
}
