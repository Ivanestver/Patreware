package config

import (
	"encoding/json"
	"os"
)

type HashesConfig struct {
	MD5HashPath    string
	SHA256HashPath string
}

type SignaturesConfig struct {
	Path string
}

type Config struct {
	Hashes     HashesConfig
	Signatures SignaturesConfig
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
