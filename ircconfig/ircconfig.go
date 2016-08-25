package ircconfig

import (
	"errors"
	"gopkg.in/yaml.v2"
)

type IrcConfig struct {
	ConnectionString string `yaml:"connection_string"`
	Nick string
	Channels []string `yaml:",flow"`
	UseSSL bool `yaml:"use_ssl"`
}

func (c *IrcConfig) Parse(data []byte) error {
	if err := yaml.Unmarshal(data, &c); err != nil {
		return err
	}

	if c.ConnectionString == "" {
		return errors.New("IrcConfig: invalid `ConnectionString`")
	}
	if c.Nick == "" {
		return errors.New("IrcConfig: invalid `Nick`")
	}

	return nil
}