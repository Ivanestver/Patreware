package config

import (
	"encoding/json"
	"os"
)

type Hashes struct {
	MD5HashPath    string
	SHA256HashPath string
}

type Config struct {
	Hashes Hashes
}

var config *Config

func InitializeConfig() {
	config = &Config{}
	content, err := os.ReadFile("config.json")
	if err != nil {
		panic(err.Error())
	}
	if err = json.Unmarshal(content, config); err != nil {
		panic(err.Error())
	}
}

func GetConfig() *Config {
	if config == nil {
		panic("The application hasn't been initialized yet")
	}
	return config
}
