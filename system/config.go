package system

import (
	"encoding/json"
	"io/ioutil"
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
	WebPath    string                 //web source path
	allValues  map[string]interface{} `json:"-"`
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
	defer r.Close()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatalln(err)
	}

	//解码JSON
	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Fatalln(err)
	}

	var allValues map[string]interface{}
	err = json.Unmarshal(data, &allValues)
	conf.allValues = allValues

	if err != nil {
		log.Fatalln(err)
	} else {
		isLoad = true
	}
}

func (this *Config) Value(key string) interface{} {
	if !isLoad {
		panic("load config file before get it!")
	}

	if _, ok := this.allValues[key]; !ok {
		return nil
	}

	return this.allValues[key]
}

func (this *Config) ToString() string {
	data, err := json.Marshal(this.allValues)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)

}
