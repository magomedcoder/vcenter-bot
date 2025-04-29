package internal

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type VCenter struct {
	Host            string `yaml:"host"`
	UsernamePostfix string `yaml:"username_postfix"`
}

type ProxySocks5 struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Telegram struct {
	Token       string       `yaml:"token"`
	ProxySocks5 *ProxySocks5 `yaml:"proxy_socks5"`
}

type Config struct {
	Telegram Telegram `yaml:"telegram"`
	VCenter  *VCenter `yaml:"vcenter"`
}

func ReadConfig(filename string) (*Config, error) {
	conf := &Config{}
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if yaml.Unmarshal(content, conf) != nil {
		panic(fmt.Sprintf("%s: %v", filename, err))
	}

	return conf, nil
}
