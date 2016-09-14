package config

import (
	"io/ioutil"
	"log"

	"github.com/burntsushi/toml"
)

// Config manages configuration for Charon
type Config struct {
	ServiceTimeout string `toml:"service_timeout"`
	Port           string
	Services       map[string]ServiceDefinition
}

// ServiceDefinition represents the URL and desired routing prefix
// for a given service.
//
//		The URL part MUST be of the form "http://hostname:port".
//
//		The prefix MUST be a string in the form of a glob, e.g. "/service/*"
//		will route all requests to the service with a prefix of "/service".
//
//		If the original request was for "http://hostname:port/users/1",
//		the routed URL will be "http://hostname:port/service/users/1".
type ServiceDefinition struct {
	URL    string
	Prefix string
}

// Load reads config data from a configuration file
func (c *Config) Load(path string, logger *log.Logger) *Config {
	tomlConfig, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := toml.Decode(string(tomlConfig), &c); err != nil {
		log.Fatal(err)
	}

	var serviceCount = len(c.Services)
	if serviceCount > 0 {
		logger.Printf("info: found %d services.", serviceCount)
	}

	return c
}
