package config

import (
	"fmt"

	"github.com/saromanov/cowrow"
)

// Config defines configuration for mystery
type Config struct {
	Server  Server  `yaml:"server"`
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

// Server defines configuration for server
type Server struct {
	Address string `yaml:"address"`
	Key     string `yaml:"key"`
	Crt     string `yaml:"crt"`
}

// Load provides loading of the config
func Load(path string) (*Config, error) {

	c := &Config{}
	if err := cowrow.LoadByPath(path, &c); err != nil {
		return makeDefault(), fmt.Errorf("unable to load config: %v", err)
	}

	return c, nil
}

func makeDefault() *Config {
	return &Config{
		Server: Server{
			Address: ":8085",
		},
	}
}
