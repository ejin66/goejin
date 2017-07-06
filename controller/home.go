package controller

import (
	"strconv"
	"GoEjin/system/controller"
	"GoEjin/db"
	"fmt"
)

type HomeController struct {
	controller.BaseController
}

//func (self *HomeController) Filter() (bool,string) {
//	return true,""
//}

/*              以下方法为自定义方法             */

func (this *HomeController) Index() {
	fmt.Printf("%p,%p",this,&this.BaseController)
	body := this.Ctx.Body()
	this.Ctx.Response("index....." + body)
}

func (this *HomeController) Show() {
	this.Ctx.Response("show data")
}

func (this *HomeController) Add(i string, j string) {
	ii, err := strconv.Atoi(i)
	ij, err2 := strconv.Atoi(j)

	if err != nil || err2 != nil {
		this.Ctx.Response("arguments type error")
	}
	this.Ctx.Response("计算结果：" + strconv.Itoa(ii+ij))
}

func (this *HomeController) Insert(name string) {
	result := db.Insert("user_info", db.Ipt{"user_name": name})
	if result == -1 {
		this.Ctx.Response("insert failed!")
		return
	}
	this.Ctx.Response("insert successful! Row number:" + strconv.Itoa(int(result)))
}
