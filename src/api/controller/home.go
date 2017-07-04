package controller

import (
	"fmt"
	"strconv"
)

type HomeController struct{
	BaseController
}

func (self  *HomeController) Instance() Base {
	return &HomeController{}
}

func (self *HomeController) Index()  {
	fmt.Printf("home:%p\n",self)
	self.Ctx.Response("index.....")
}

func (self *HomeController) Show() {
	self.Ctx.Response("show data")
}

func (self *HomeController) Add(i string,j string) {
	ii,err := strconv.Atoi(i)
	ij,err2 := strconv.Atoi(j)

	if err != nil || err2 != nil {
		self.Ctx.Response("arguments type error")
	}
	self.Ctx.Response("计算结果："+strconv.Itoa(ii + ij))
}