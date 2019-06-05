package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	IP_ADDRESS         string
	IP_PORT            string
	DB_USER            string
	DB_PASSWORD        string
	DB_ADDRESS         string
	DB_PORT            string
	DB_NAME            string
	DEFAULT_CONTROLLER string
	IMAGE_PATH         string
	DOWNLOAD_PATH      string
	BASE_URL           string
}

//生成一个全局的conf变量存储读取的配置
var conf Config

func init() {
	LoadConf()
}

func GetConfig() *Config {
	return &conf
}

//读取配置函数
func LoadConf() {
	//打开文件
	r, err := os.Open("src/GoEjin/config/config.json")
	if err != nil {
		log.Fatalln(err)
	}
	//解码JSON
	decoder := json.NewDecoder(r)
	err = decoder.Decode(&conf)
	if err != nil {
		log.Fatalln(err)
	}
}
