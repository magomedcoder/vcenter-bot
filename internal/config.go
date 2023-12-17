package internal

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type vCenter struct {
	Host     string `json:"host" yaml:"host"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type Config struct {
	TelegramToken string  `json:"telegram_token" yaml:"telegram_token"`
	Vcenter       vCenter `json:"vcenter" yaml:"vcenter"`
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

const tokenFileName = "/tmp/vcenter-bot-token"

func WriteTokenToFile(token string) {
	err := ioutil.WriteFile(tokenFileName, []byte(token), 0644)
	if err != nil {
		fmt.Errorf("ошибка записи токена в файл: %v", err)
	}
}

func readTokenFromFile() (string, error) {
	data, err := ioutil.ReadFile(tokenFileName)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения токена из файла: %v", err)
	}
	return string(data), nil
}
