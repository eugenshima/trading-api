// Package config provides configuration information
package config

import "github.com/caarlos0/env"

type Config struct {
	SigningKey string `env:"SIGNING_KEY" envDefault:"ew4t137tr1eyfg1ryg4ryerg2743gr2"`
}

// NewConfig creates a new Config instance
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
