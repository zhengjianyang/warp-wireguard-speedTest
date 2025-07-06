package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Threads      int      `yaml:"threads"`
	PingCount    int      `yaml:"pingCount"`
	Ports        []int    `yaml:"ports"`
	IPv4CIDRs    []string `yaml:"ipv4CIDRs"`
	SaveFileName string   `yaml:"saveFileName"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
