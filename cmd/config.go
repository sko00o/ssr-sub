package main

import (
	"io/ioutil"

	"github.com/mingcheng/ssr-subscriber/node"
	"gopkg.in/yaml.v3"
)

// Configure struct, for more information see config-example.yaml file
type Configure struct {
	URL       []string         `yaml:"url"`
	File      []string         `yaml:"file"`
	Proxy     string           `yaml:"proxy"`
	Check     node.CheckConfig `yaml:"check"`
	Interval  int              `yaml:"interval"`
	Bind      string           `yaml:"bind"`
	RedisAddr string           `yaml:"redis"`
}

// ParseConfig to parse config via local filepath
func ParseConfig(configPath string) (*Configure, error) {
	var configure Configure

	// read configure from config file
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &configure)
	if err != nil {
		return nil, err
	}

	return &configure, err
}
