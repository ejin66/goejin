//@author jsl
//@time 2019-06-05 17:15
package goejin

import (
	"fmt"
	"github.com/ejin66/goejin/system"
	"github.com/ejin66/goejin/util"
	"net/http"
)

func Listen(path string, router system.Router) {
	system.LoadConf(path)
	system.LoadRouter(router)
	ip := system.GetConfig().IpAddress + ":" + system.GetConfig().IpPort

	fmt.Println("server listen on1 ", ip)
	err := http.ListenAndServe(ip, system.GetServeMux())
	if err != nil {
		util.PrintError("server listen err:" + err.Error())
	} else {
		fmt.Println("server listen on ", ip)
	}
}
