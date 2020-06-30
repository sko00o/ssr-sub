package main

import subscriber "github.com/mingcheng/ssr-subscriber"

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
