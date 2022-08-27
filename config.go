package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Telegram struct {
		ApiToken string `yaml:"api_token"`
	}
	CouchDB struct {
		Host     string
		User     string
		Password string
		Database string
	}
}

func getConfig(filename string) (config *Config, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	config = new(Config)
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return
	}
	return
}
