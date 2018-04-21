package util

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/qiniu/log"
)

type configuration struct {
	RedisDB RedisDBConfiguration `yaml:"redis"`
}

// RedisDBConfiguration redis
type RedisDBConfiguration struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}

var Configuration *configuration

//func initConfiguration(configFile string) {
//	loadConfiguration(configFile)
//}

func initConfig() {
	loadConfiguration()
}

func loadConfiguration() {
	configFilePath := "util/config.yaml"
	//configFilePath := fmt.Sprintf("common/%s", configFile)
	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("[loadConfiguration] Relative path error %v \n", err)
	}
	err = yaml.Unmarshal(file, &Configuration)
	if err != nil {
		log.Fatalf("[loadConfiguration]:%v\n", err)
	}
}
