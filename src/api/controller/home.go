package controller

import (
	"strconv"
	"api/system"
)

/*              以下结构或方法为必要             */
type HomeController struct {
	system.BaseController
}

func (self *HomeController) Instance() system.Base {
	return &HomeController{}
}

/*              以下方法为自定义方法             */

func (self *HomeController) Index() {
	body := self.Ctx.Body()
	self.Ctx.Response("index....." + body)
}

func (self *HomeController) Show() {
	self.Ctx.Response("show data")
}

func (self *HomeController) Add(i string, j string) {
	ii, err := strconv.Atoi(i)
	ij, err2 := strconv.Atoi(j)

	if err != nil || err2 != nil {
		self.Ctx.Response("arguments type error")
	}
	self.Ctx.Response("计算结果：" + strconv.Itoa(ii+ij))
}
