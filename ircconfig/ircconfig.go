package ircconfig

import (
	"errors"
	"math/rand"
	"gopkg.in/yaml.v2"
)

type IrcConfig struct {
	ConnectionString string `yaml:"connection_string"`
	Nick string
	Channels []string `yaml:",flow"`
	UseSSL bool `yaml:"use_ssl"`
	Greetings []string `yaml:",flow"`
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

	if len(c.Greetings) == 0 {
		c.Greetings = []string{"hi"} //put something in there so it can still talk
	}

	return nil
}

func (c IrcConfig) GetRandomGreeting() string {
	return c.Greetings[rand.Intn(len(c.Greetings))]
}