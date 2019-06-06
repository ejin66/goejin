package system

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	IpAddress  string
	IpPort     string
	DbUser     string
	DbPassword string
	DbAddress  string
	DbPort     string
	DbName     string
	WebPath    string //web source path
}

//生成一个全局的conf变量存储读取的配置
var conf Config

var isLoad = false

func GetConfig() *Config {
	if !isLoad {
		panic("load config file before get it!")
	}

	return &conf
}

//读取配置函数
func LoadConf(path string) {
	//打开文件
	r, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	//解码JSON
	decoder := json.NewDecoder(r)
	err = decoder.Decode(&conf)
	if err != nil {
		log.Fatalln(err)
	} else {
		isLoad = true
	}
}

func (this *Config) ToString() string {
	return "IpAddress: " + this.IpAddress +
		", IpPort: " + this.IpPort +
		", DbUser: " + this.DbUser +
		", DbPassword: " + this.DbPassword +
		", DbAddress: " + this.DbAddress +
		", DbPort: " + this.DbPort +
		", DbName: " + this.DbName

}
