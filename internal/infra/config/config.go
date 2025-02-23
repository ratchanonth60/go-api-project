package config

import (
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v11"
	"gopkg.in/yaml.v3"
)

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
