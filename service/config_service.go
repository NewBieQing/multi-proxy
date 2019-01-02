package service

import (
	"io/ioutil"
	"github.com/BurntSushi/toml"
)

//设置系统环境变量 dev docker online
var ENV = "dev"
var ConfigPath = ""

type Config struct {
	WebLog WebLogConfig
}

type WebLogConfig struct {
	FilePath string
	Level    string
}

type ProxysConfig struct {
	Proxy []ProxyConfig
}

type ProxyConfig struct {
	From string
	To   string
}

var (
	config *Config
)

func LoadWebLogConfig() *WebLogConfig {
	return &LoadConfig().WebLog
}


func LoadProxysConfig() *ProxysConfig{
	proxys := &ProxysConfig{}
	data, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		panic(err)
	}
	if _, err := toml.Decode(string(data), &proxys); err != nil {
		panic(err)
	}
	return proxys
}

func LoadConfig() *Config {
	if (config == nil) {
		config = &Config{}
		// read the raw contents of the file
		data, err := ioutil.ReadFile(ConfigPath)
		if err != nil {
			panic(err)
		}
		// put the file's contents as toml to the default configuration(c)
		if _, err := toml.Decode(string(data), &config); err != nil {
			panic(err)
		}
	}
	return config
}

func SetConfigPath(path string) {
	ConfigPath = path
}
