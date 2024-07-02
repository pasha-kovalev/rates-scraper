package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

const configFilename = "config.yaml"

type Config struct {
	App struct {
		LogLevel string `yaml:"log-level"`
		Version  string
		Host     string
		Port     int
	}
	Db struct {
		User     string
		Password string
		Host     string
		Port     int
		Name     string
	}
}

func LoadConfig() (Config, error) {
	confBytes, err := os.ReadFile(configFilename)
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = yaml.Unmarshal(confBytes, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
