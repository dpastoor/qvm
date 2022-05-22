package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DefaultVersion string `yaml:"default_version"`
}

func (cfg Config) Validate() error {
	return nil
}

func read(path string) (Config, error) {
	var config Config
	file, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}
	err = config.Validate()
	if err != nil {
		return config, err
	}
	return config, nil
}
