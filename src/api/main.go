package main

import (
	"net/http"
	"fmt"
	"api/config"
	"api/system"
)

var ip string
var serverMux *http.ServeMux

func init() {
	ip = config.IP_ADDRESS + ":" + config.IP_PORT
	serverMux = system.GetServeMux()
}

func main() {
	fmt.Println("Server listen on ", ip)
	err := http.ListenAndServe(ip ,serverMux)
	if err != nil {
		fmt.Println("server listen err:",err.Error())
	}
}

