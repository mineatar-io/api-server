package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Redis struct {
		URI      string `yaml:"uri"`
		Database int    `yaml:"database"`
	} `yaml:"redis"`
}

func (c *Configuration) ReadFile(file string) error {
	data, err := ioutil.ReadFile(file)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, c)
}
