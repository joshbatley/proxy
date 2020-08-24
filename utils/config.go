package utils

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var LoadedConfig Config

// Config file structure
type Config struct {
	Env    string `yml:"env"`
	Name   string `yml:"name"`
	Port   string `yml:"port"`
	DBFile string `yml:"DBFile"`
}

// LoadConfig read from file location
func LoadConfig(f string) (*Config, error) {
	file, _ := ioutil.ReadFile(f)
	config := &Config{}
	err := yaml.Unmarshal([]byte(file), &config)

	if err != nil {
		return nil, err
	}
	if config != nil {
		LoadedConfig = *config
	}
	return config, nil
}
