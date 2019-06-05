package main

import (
	sysController "github.com/ejin66/goejin/system"
)

/**
Route table
POST/GET/EMPTY STRING(both post and get)
*/
var RouteTable = sysController.Router{
	"HOME": sysController.Cfg{&HomeController{}, "", sysController.MethodMap{"Index": "GET"}},
}
