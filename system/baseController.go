package system

import (
	"encoding/json"
	"errors"
	"github.com/ejin66/goejin/system/session"
	"github.com/ejin66/goejin/util"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Base interface {
	Context(c *Context)     //传递请求上下文
	Filter() (bool, string) //该方法在接口方法调用之前调用，用来过滤部分请求
}

/**
Cb : controller
DefaultMethod : default controller request method. if "",means no limit for method
Methods : controller functions request method. if "",means no limit for method
*/
type Cfg struct {
	Cb            Base
	DefaultMethod string
	Methods       MethodMap
}

type MethodMap map[string]string

type Router map[string]Cfg

type Context struct {
	W        *http.ResponseWriter
	Req      *http.Request
	postJson *map[string]interface{}
}

func (this *Context) Response(msg string) {
	io.WriteString(*(this.W), msg)
}

func (this *Context) AddHeader(key, value string) {
	(*this.W).Header().Add(key, value)
}

func (this *Context) setStatusCode(code int) {
	(*this.W).WriteHeader(code)
}

/*
获取post参数
*/
func (this *Context) FetchForm(key string) string {
	util.Print(key, this.Req.PostFormValue(key))
	return this.Req.PostFormValue(key)
}

func (this *Context) FetchBodyAsJson() string {
	body, _ := ioutil.ReadAll(this.Req.Body)
	return string(body)
}

/**
pointer
*/
func (this *Context) FetchBodyAs(i interface{}) {
	body, _ := ioutil.ReadAll(this.Req.Body)
	json.Unmarshal(body, i)
}

func (this *Context) FetchJson(key string) (interface{}, error) {
	if len(*this.postJson) == 0 {
		body, _ := ioutil.ReadAll(this.Req.Body)
		json.Unmarshal(body, this.postJson)
	}
	if _, ok := (*this.postJson)[key]; !ok {
		return nil, errors.New("can't find key " + key)
	}
	return (*this.postJson)[key], nil
}

func (this *Context) FetchJsonInt(key string) (returnValue int, returnError error) {
	defer func() {
		if err := recover(); err != nil {
			util.PrintError(err)
			returnError = errors.New("value type is not int")
		}
	}()
	v, err := this.FetchJson(key)
	if err == nil {
		returnValue = int(v.(float64))
	}
	returnError = err
	return
}

func (this *Context) FetchJsonString(key string) (returnValue string, returnError error) {
	defer func() {
		if err := recover(); err != nil {
			returnError = errors.New("value type is not string")
		}
	}()
	v, err := this.FetchJson(key)
	if err == nil {
		returnValue = v.(string)
	}
	returnError = err
	return
}

func (this *Context) FetchFileAs(key string, path string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			err = errors.New("no file found: " + key)
		}
	}()

	file, _, err := this.Req.FormFile(key)

	if err != nil {
		return
	}

	lf, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer lf.Close()

	_, err = io.Copy(lf, file)
	return
}

/*
获取post body
*/
func (this *Context) Body() string {
	body, err := ioutil.ReadAll(this.Req.Body)

	if err != nil {
		util.PrintError(err)
		return ""
	}
	return string(body)
}

func (this *Context) SessionStart() session.Session {
	return session.SessionManager.SessionStart(this.W, this.Req)
}

type BaseController struct {
	Ctx *Context
}

func (this *BaseController) Context(c *Context) {
	this.Ctx = c
	json := make(map[string]interface{})
	c.postJson = &json
}

func (this *BaseController) Filter() (bool, string) {
	return true, ""
}
