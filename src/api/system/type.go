package system
/**
Cb : controller
Method : default controller request method
Methods : controller functions exact request method
 */
type Cfg struct {
	Cb Base
	DefaultMethod string
	Methods MethodMap
}

type MethodMap map[string]string

type Router map[string]Cfg
