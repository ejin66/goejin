package controller

import (
	"fmt"
	"strconv"
	"api/system"
)

/*              以下结构或方法为必要             */
type HomeController struct{
	system.BaseController
}

func (self  *HomeController) Instance() system.Base {
	return &HomeController{}
}

/*              以下方法为自定义方法             */

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