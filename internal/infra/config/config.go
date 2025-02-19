package config

import (
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v11"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Server struct {
		Port string `yaml:"port" env:"PORT"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host" env:"POSTGRES_HOST"`
		Port     string `yaml:"port" env:"POSTGRES_PORT"`
		User     string `yaml:"user" env:"POSTGRES_USER"`
		Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
		DBName   string `yaml:"db_name" env:"POSTGRES_DB"`
	} `yaml:"database"`
	Cache struct {
		Address     string `yaml:"address" env:"REDIS_HOST"`
		Port        string `yaml:"port" env:"REDIS_PORT"`
		MaxSizePool string `yaml:"max_size_pool" env:"REDIS_MAX_POOL"`
	} `yaml:"cache"`
	JWT struct {
		Signed string `yaml:"signed" env:"JWT_SIGNED"`
	} `yaml:"jwt"`
}

var (
	Config *AppConfig
	IsYaml = true
)

func LoadConfig(path string) error {
	Config = &AppConfig{}
	// load config from yaml file
	if IsYaml {
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read config file: %v", err)
		}
		if err := yaml.Unmarshal(data, Config); err != nil {
			log.Fatalf("failed to unmarshal config file: %v", err)
		}

	} else {
		if err := env.Parse(Config); err != nil {
			fmt.Println(err)
			return fmt.Errorf("failed to load env: %v", err)
		}

	}
	return nil
}
