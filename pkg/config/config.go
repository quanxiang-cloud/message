package config

import (
	"io/ioutil"

	"git.internal.yunify.com/qxp/misc/mysql2"
	"git.internal.yunify.com/qxp/misc/redis2"
	"github.com/quanxiang-cloud/message/pkg/client"

	"gopkg.in/yaml.v2"
)

// Conf 配置文件
var Conf *Config

// DefaultPath 默认配置路径
var DefaultPath = "./configs/config.yml"

// Config 配置文件
type Config struct {
	Port        string        `yaml:"port"`
	Model       string        `yaml:"model"`
	InternalNet client.Config `yaml:"internalNet"`
	Mysql       mysql2.Config `yaml:"mysql"`
	Redis       redis2.Config `yaml:"redis"`
}

// NewConfig 获取配置配置
func NewConfig(path string) (*Config, error) {
	if path == "" {
		path = DefaultPath
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, &Conf)
	if err != nil {
		return nil, err
	}

	return Conf, nil
}
