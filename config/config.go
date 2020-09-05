package config

// Config defines configuration for mystery
type Config struct {
	Backend Backend
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
