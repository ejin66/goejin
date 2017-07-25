package main

import (
	"net/http"
	"fmt"
	"GoEjin/system/router"
	_ "GoEjin/system/session"
	"GoEjin/system/config"
)

var ip string

func init() {
	ip = config.GetConfig().IP_ADDRESS + ":" + config.GetConfig().IP_PORT
}

func main() {
	fmt.Println("Server listen on ", ip)
	err := http.ListenAndServe(ip ,router.GetServeMux())
	if err != nil {
		fmt.Println("server listen err:",err.Error())
	}
}

