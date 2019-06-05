package system

import (
	"bytes"
	"fmt"
	"github.com/ejin66/goejin/util"
	"os"
)

func web(path string) string {
	file, err := os.OpenFile(GetConfig().WebPath+"/"+path, os.O_RDONLY, os.ModePerm)
	defer file.Close()
	if err != nil {
		fmt.Println("open web page err:", err.Error())
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
			fmt.Println("read web page err:", err.Error(), i)
			break
		}
		buf.Write(buffer[:i])
	}
	return buf.String()
}
