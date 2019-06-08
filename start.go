//@author jsl
//@time 2019-06-05 17:15
package goejin

import (
	"github.com/ejin66/goejin/system"
	"github.com/ejin66/goejin/util"
	"net/http"
	"time"
)

func Listen(path string, router system.Router) {
	system.LoadConf(path)
	system.LoadRouter(router)
	ip := system.GetConfig().IpAddress + ":" + system.GetConfig().IpPort

	go func() {
		time.Sleep(time.Second)
		util.PrintLogDivider()
		defer util.PrintLogDivider()
		util.Print("Running successful!")
		util.Print("Listened on:", ip)
	}()

	err := http.ListenAndServe(ip, system.GetServeMux())
	if err != nil {
		util.PrintError(err.Error())
	}
}
