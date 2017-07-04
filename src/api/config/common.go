package config

import (
	"os"
	"fmt"
	"bytes"
)

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
