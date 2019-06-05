package main

import "goejin"

func main() {
	goejin.Listen("./config.json", RouteTable)
}
