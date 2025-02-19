package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"db_name"`
	} `yaml:"database"`
	Cache struct {
		Address     string `yaml:"address"`
		Port        string `yaml:"port"`
		MaxSizePool int    `yaml:"max_size_pool"`
	} `yaml:"cache"`
	JWT struct {
		Signed string `yaml:"signed"`
	} `yaml:"jwt"`
}

var Config *AppConfig

func LoadConfig(path string) error {
	// load config from yaml file
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}
	if err := yaml.Unmarshal(data, &Config); err != nil {
		log.Fatalf("failed to unmarshal config file: %v", err)
	}
	return nil
}
