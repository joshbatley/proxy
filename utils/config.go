package utils

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config - config file structure
type Config struct {
	Name   string `yml:"name"`
	Port   string `yml:"port"`
	DBFile string `yml:"DBFile"`
}

// LoadConfig - Read from file location
func LoadConfig(f string) (*Config, error) {
	file, _ := ioutil.ReadFile(f)
	config := &Config{}
	err := yaml.Unmarshal([]byte(file), &config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
