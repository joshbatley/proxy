package utils

import (
	"io/ioutil"
	"net/url"
	"regexp"

	"gopkg.in/yaml.v2"
)

// FormatURL - find url in query param
func FormatURL(u string) *url.URL {
	s := regexp.MustCompile(`(?:/query\?q=)(.{0,})`)
	r := string(s.ReplaceAll([]byte(u), []byte("$1")))
	formattedURL, err := url.Parse(r)
	if err != nil {
		panic(err)
	}
	return formattedURL
}

// Config - config file structure
type Config struct {
	Name string `yml:"name"`
	Port string `yml:"port"`
}

// ReadConfig - Read from file location
func ReadConfig(f string) (*Config, error) {
	file, _ := ioutil.ReadFile(f)
	config := &Config{}
	err := yaml.Unmarshal([]byte(file), &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
