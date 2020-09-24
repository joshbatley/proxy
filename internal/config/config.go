package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config file structure
type Config struct {
	Env    string `yml:"env"`
	Name   string `yml:"name"`
	Port   string `yml:"port"`
	DBFile string `yml:"DBFile"`
}

// Load read from file location
func Load(f string) (*Config, error) {
	file, _ := ioutil.ReadFile(f)
	config := &Config{}
	err := yaml.Unmarshal([]byte(file), &config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
