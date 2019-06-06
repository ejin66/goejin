package system

import (
	"bytes"
	"github.com/ejin66/goejin/util"
	"os"
	"strings"
)

func web(path string) string {
	sourcePath := GetConfig().WebPath + "/" + path
	if strings.HasSuffix(sourcePath, "/") {
		sourcePath += "index.html"
	}
	util.Print("load static resource: ", sourcePath)
	file, err := os.OpenFile(sourcePath, os.O_RDONLY, os.ModePerm)
	defer file.Close()
	if err != nil {
		util.PrintError("open static resource,", err.Error())
		return util.Error404Web()
	}
	buf := bytes.NewBufferString("")
	for {
		var buffer = make([]byte, 1024*10)
		i, err := file.Read(buffer)
		if i <= 0 {
			break
		}
		if err != nil {
			util.PrintError("read static resource,", err.Error(), i)
			break
		}
		buf.Write(buffer[:i])
	}
	return buf.String()
}
