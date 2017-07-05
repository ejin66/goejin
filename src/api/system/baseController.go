package system

import (
	"net/http"
	"io"
)

type Base interface {
	Context(c *Context)
	Instance() Base //确保每次请求，都开辟新的地址空间，避免多个请求共用一个
}

type Context struct {
	W *http.ResponseWriter
	Req *http.Request
}

func (this *Context) Response(msg string) {
	io.WriteString(*(this.W), msg)
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

