package config

import (
	"fmt"

	"github.com/saromanov/cowrow"
)

// Config defines configuration for mystery
type Config struct {
	Backend Backend `yaml:"backend"`
	// Environment variable for master pass
	MasterPassEnv string `yaml:"masterPassEnv"`
}

// Backend defines way for store secrets
type Backend struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	DB       string `yaml:"db"`
	Password string `yaml:"password"`
}

// Load provides loading of the config
func Load(path string) (*Config, error) {

	c := &Config{}
	if err := cowrow.LoadByPath(path, &c); err != nil {
		return nil, fmt.Errorf("unable to load config: %v", err)
	}

	return c, nil
}
