package config

import (
	"io/ioutil"
	"time"

	"git.internal.yunify.com/qxp/misc/client"
	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/misc/mysql2"

	"gopkg.in/yaml.v2"
)

// Conf 配置文件
var Conf *Config

// DefaultPath 默认配置路径
var DefaultPath = "./configs/config.yml"

// Config 配置文件
type Config struct {
	Port         string        `yaml:"port"`
	Model        string        `yaml:"model"`
	MessageAPI   string        `yaml:"messageAPI"`
	InternalNet  client.Config `yaml:"internalNet"`
	ProcessorNum int           `yaml:"processorNum"`
	SyncChannel  string        `yaml:"syncChannel"`
	HandOut      HandOut       `yaml:"handout"`
	Log          logger.Config `yaml:"log"`
	Mysql        mysql2.Config `yaml:"mysql"`
	Email        Email         `yaml:"email"`
	AUth         Auth          `yaml:"auth"`
}

// Email email
type Email struct {
	EmailList []EmailConfig `yaml:"emails"`
}

// EmailConfig list
type EmailConfig struct {
	Emailfrom string `yaml:"emailfrom"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Aliasname string `yaml:"aliasname"`
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

// HandOut hand out
type HandOut struct {
	Deadline     time.Duration
	DialTimeout  time.Duration
	MaxIdleConns int
}

//Auth token check
type Auth struct {
	CheckToken string `yaml:"checktoken"`
}
