package config

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type database struct {
	Address  string
	Port     int
	Dbname   string
	User     string
	Password string
}

type app struct {
	Port    string
	Mode    string
	Timeout int
	Secret  string
}

type t struct {
	Database database
	App      app
}

var Config t

var cliPath = flag.String("path", "./dev.yaml", "Input config file path")

func (config *t) GetConfig() {
	// 得到参数
	flag.Parse()
	if *cliPath == "" {
		panic("Please input correct config file path")
	}

	// 打开文件
	file, err := ioutil.ReadFile(*cliPath)
	if err != nil {
		panic("Open config file error")
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		panic(err.Error())
	}
}

