package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigServer struct {
	Addr string `yaml:"addr"`
}

type Config struct {
	Server      ConfigServer `yaml:"server"`
	ProxyServer ConfigServer `yaml:"proxy_server"`
}

var SysConfig *Config = &Config{}

func init() {
	path := "/home/ogios/.config/transfer-go/settings.yml"
	b, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Config file fail to read: %s", err.Error()))
	}
	err = yaml.Unmarshal(b, SysConfig)
	if err != nil {
		panic("Config fail to load")
	}
}
