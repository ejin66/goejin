package controller

import (
	"net/http"
	"io"
	"fmt"
	"io/ioutil"
	"api/system/common"
)

type Base interface {
	Context(c *Context)
	Instance() Base //确保每次请求，都开辟新的地址空间，避免多个请求共用一个
}

/**
Cb : controller
Method : default controller request method. if "",means no limit for method
Methods : controller functions exact request method. if "",means no limit for method
 */
type Cfg struct {
	Cb Base
	DefaultMethod string
	Methods MethodMap
}

type MethodMap map[string]string

type Router map[string]Cfg

type Context struct {
	W   *http.ResponseWriter
	Req *http.Request
}

func (this *Context) Response(msg string) {
	io.WriteString(*(this.W), msg)
}

/*
获取post参数
 */
func (this *Context) Fetch(key string) string {
	fmt.Println(key, this.Req.PostFormValue(key))
	//body, _ := ioutil.ReadAll(this.Req.Body)
	return this.Req.PostFormValue(key)
}

/*
获取post body
 */
func (this *Context) Body() string {
	body, err := ioutil.ReadAll(this.Req.Body)

	if err != nil {
		common.PrintError(err)
		return ""
	}
	return string(body)
}

type BaseController struct {
	Ctx *Context
}

func (this *BaseController) Instance() Base {
	return &BaseController{}
}

func (this *BaseController) Context(c *Context) {
	this.Ctx = c
}
