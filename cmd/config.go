package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mingcheng/ssr-subscriber"
	"gopkg.in/yaml.v3"
)

// Configure struct, for more information see config-example.yaml file
type Configure struct {
	URL       []string               `yaml:"url"`
	File      []string               `yaml:"file"`
	Output    string                 `yaml:"output"`
	Proxy     string                 `yaml:"proxy"`
	Check     subscriber.CheckConfig `yaml:"check"`
	Interval  int                    `yaml:"update_interval"`
	Bind      string                 `yaml:"bind"`
	AutoClean bool                   `yaml:"auto_clean"`
	Exceed    uint                   `yaml:"config_exceed"`
}

// ParseConfig to parse config via local filepath
func ParseConfig(configPath string) (*Configure, error) {
	var configure Configure

	// read configure from config file
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &configure); err != nil {
		return nil, err
	}

	// check output directory
	if stat, err := os.Stat(configure.Output); err != nil || !stat.IsDir() {
		return nil, fmt.Errorf("%s is not directory or not wirtable", configure.Output)
	}

	return &configure, err
}
