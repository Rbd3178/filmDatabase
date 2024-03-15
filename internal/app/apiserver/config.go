package apiserver

// Config
type Config struct {
	Port     string `json:"port"`
	LogLevel string `json:"log_level"`
}

// NewConfig
func NewConfig() *Config {
	return &Config{
		Port:     "8080",
		LogLevel: "debug",
	}
}
