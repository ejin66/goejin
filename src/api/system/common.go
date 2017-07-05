package system

import (
	"os"
	"fmt"
	"bytes"
)

func PrintError(msg interface{}) {
	//其中0x1B是标记，[开始定义颜色，1代表高亮，39代表黑色背景，31代表红色前景，0代表恢复默认颜色。
	fmt.Printf("%c[0;39;31m%s%c[0m\n", 0x1B ,msg, 0x1B)
}

func Error404() string {
	return "404 not found"
}

func Web(path string) string {
	file, err := os.OpenFile("web/"+path, os.O_RDONLY, os.ModePerm)
	defer file.Close()
	if err != nil {
		fmt.Println("open web page err:", err.Error())
		return Error404()
	}
	buf := bytes.NewBufferString("")
	for {
		var buffer []byte = make([]byte, 1024*10)
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
