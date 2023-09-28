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

var GLOBAL_CONFIG *Config = &Config{}
var CONFIG_FILE_PATH string = "/transfer-go/base_server.yml"

func get_path_with_args() (string, bool) {
	args := os.Args
	for index, arg := range args {
		if arg == "-c" {
			if len(args) <= index+1 {
				panic(fmt.Sprintln("No path provided for -c"))
			} else {
				return args[index+1], true
			}
		}
	}
	return "", false
}

func init() {
	path, ok := get_path_with_args()
	if !ok {
		defalut_path, err := os.UserConfigDir()
		if err != nil {
			panic(fmt.Sprintf("Config directory not found: %s", err.Error()))
		}
		path = defalut_path + CONFIG_FILE_PATH
	}
	b, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Config file fail to read: %s", err.Error()))
	}
	err = yaml.Unmarshal(b, GLOBAL_CONFIG)
	if err != nil {
		panic("Config fail to load")
	}
}
