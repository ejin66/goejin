package system
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
