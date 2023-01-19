package config

import (
	"github.com/kelseyhightower/envconfig"
)

// ReadEnv reads some configs from environment variables
func ReadEnv(cfg *Config) error {
	return envconfig.Process("", cfg)
}
