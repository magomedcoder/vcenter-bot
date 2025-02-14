package internal

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type VCenter struct {
	Host            string `json:"host" yaml:"host"`
	UsernamePostfix string `json:"username_postfix" yaml:"username_postfix"`
}

type Config struct {
	TelegramToken string   `json:"telegram_token" yaml:"telegram_token"`
	VCenter       *VCenter `json:"vcenter" yaml:"vcenter"`
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
