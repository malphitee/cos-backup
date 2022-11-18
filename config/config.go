package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Url               string `yaml:"url" mapstructure:"url"`
	SecretId          string `yaml:"secret_id" mapstructure:"secret_id"`
	SecretKey         string `yaml:"secret_key" mapstructure:"secret_key"`
	Dir               string `yaml:"dir" mapstructure:"dir"`
	Path              string `yaml:"path" mapstructure:"path"`
	DeleteRemote      bool   `yaml:"delete_remote" mapstructure:"delete_remote"`
	ServerChanSendKey string `yaml:"serverchan_send_key" mapstructure:"serverchan_send_key"`
}

func GetConfigFromYaml() Config {
	var currentPath string
	var err error
	if currentPath, err = os.Getwd(); err != nil {
		panic("获取当前目录出错，err = " + err.Error())
	}
	vip := viper.New()
	vip.AddConfigPath(currentPath)
	vip.SetConfigName("config")
	vip.SetConfigType("yaml")
	if err = vip.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("读取配置文件出错，err = %v", err))
	}
	var config Config
	if err = vip.Unmarshal(&config); err != nil {
		panic(fmt.Sprintf("解析配置文件出错，err = %v", err))
	}
	return config
}
