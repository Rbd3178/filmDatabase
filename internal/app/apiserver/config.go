package apiserver

import "github.com/Rbd3178/filmDatabase/internal/app/store"

// Config
type Config struct {
	Port     string `toml:"port"`
	LogLevel string `toml:"log_level"`
	Store *store.Config
}

// NewConfig
func NewConfig() *Config {
	return &Config{
		Port:     "8080",
		LogLevel: "debug",
		Store: store.NewConfig(),
	}
}
