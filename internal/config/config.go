package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Port      string `env:"PORT" envDefault:"8080"`
	DBURL     string `env:"DB_URL" envDefault:"postgres://armpower:armpower@localhost:5432/armpower?sslmode=disable"`
	JWTSecret string `env:"JWT_SECRET" envDefault:"arm-power-secret-key-change-in-production"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
