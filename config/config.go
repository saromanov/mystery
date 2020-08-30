package config

// Config defines configuration for mystery
type Config struct {
	MasterPassBackend MasterPassBackend
}

type MasterPassBackend struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	DB       string `yaml:"db"`
	Password string `yaml:"password"`
}
