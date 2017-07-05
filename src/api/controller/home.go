package controller

import (
	"strconv"
	"api/system/controller"
	"api/db"
)

/*              以下结构或方法为必要             */
type HomeController struct {
	controller.BaseController
}

func (self *HomeController) Instance() controller.Base {
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

func (self *HomeController) Insert(name string) {
	result := db.Insert("user_table", db.Ipt{"user_name" : name})
	if result == -1 {
		self.Ctx.Response("insert failed!")
		return
	}
	self.Ctx.Response("insert successful! Row number:" + strconv.Itoa(int(result)))
}
