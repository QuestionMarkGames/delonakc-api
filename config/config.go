package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

const RouteVarsKey = "VarsKey"
const (
	CommonError = -1
	NotLogin = -2
	Sucess = 0
)

type Conf struct {
	Mode string `yaml:"mode"`
	Server struct {
		Port int `yaml:"port"`
		Addr string `yaml:"addr""`
	}
	DataBase struct {
		Host string
		Port int
		Db string
	}
	Redis struct {
		Port int
		Host string
		ExpireTime string `yaml:"expireTime"`
	}
	Github
}

func Get(mode string) *Conf {

	var fileName string

	if mode == "prod" {
		fileName = "config.yml"
	} else {
		fileName = fmt.Sprintf("config-%s.yml", mode)
	}

	fmt.Println("mod environment: ", mode)

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal("Read Config Error: ", err)
	}

	var conf Conf

	err = yaml.Unmarshal(data, &conf)

	if err != nil {
		log.Fatal("Config Parse error: ", err)
	}

	return &conf
}

func (c *Conf) ToString() string {
	return fmt.Sprintf("mode: %s \nserver: %s:%d", c.Mode, c.Server.Addr, c.Server.Port)
}