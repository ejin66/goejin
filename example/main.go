package main

import "github.com/ejin66/goejin"

func main() {
	goejin.Listen("./config.json", RouteTable)
}
