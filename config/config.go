package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Url               string `yaml:"url" mapstructure:"url"`
	SecretId          string `yaml:"secret_id" mapstructure:"secret_id"`
	SecretKey         string `yaml:"secret_key" mapstructure:"secret_key"`
	Path              string `yaml:"path" mapstructure:"path"`
	DeleteRemote      bool   `yaml:"delete_remote" mapstructure:"delete_remote"`
	ServerChanSendKey string `yaml:"serverchan_send_key" mapstructure:"serverchan_send_key"`
	NotifyDriver      string `yaml:"notify_driver" mapstructure:"notify_driver"`
	GotifyToken       string `yaml:"gotify_token" mapstructure:"gotify_token"`
	GotifyUrl         string `yaml:"gotify_url" mapstructure:"gotify_url"`
}

var config *Config

func GetConfigFromYaml() {
	var currentPath string
	var err error
	if currentPath, err = os.Getwd(); err != nil {
		panic("获取当前目录出错,err = " + err.Error())
	}
	vip := viper.New()
	vip.AddConfigPath(currentPath)
	vip.SetConfigName("config")
	vip.SetConfigType("yaml")
	if err = vip.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("读取配置文件出错,err = %v", err))
	}
	if err = vip.Unmarshal(&config); err != nil {
		panic(fmt.Sprintf("解析配置文件出错,err = %v", err))
	}
}

func GetConfig() *Config {
	if config == nil {
		GetConfigFromYaml()
	}
	return config
}
