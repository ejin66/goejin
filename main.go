package main

import (
	"net/http"
	"fmt"
	"GoEjin/config"
	"GoEjin/system/router"
	_ "GoEjin/system/session"
)

var ip string

func init() {
	ip = config.IP_ADDRESS + ":" + config.IP_PORT
}

func main() {
	fmt.Println("Server listen on ", ip)
	err := http.ListenAndServe(ip ,router.GetServeMux())
	if err != nil {
		fmt.Println("server listen err:",err.Error())
	}
}

