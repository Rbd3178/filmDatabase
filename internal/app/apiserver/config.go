package apiserver

// Config
type Config struct {
	Port     string `toml:"port"`
	LogLevel string `toml:"log_level"`
	DatabaseURL string `toml:"database_url"`
}

// NewConfig
func NewConfig() *Config {
	return &Config{
		Port:     "8080",
		LogLevel: "debug",
	}
}
