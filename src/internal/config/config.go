package config

import (
	"github.com/caarlos0/env/v10"
	"log"
	"net"
)

type Config struct {
	Host        string `env:"SERVER_HOST" env-default:"localhost"`
	Port        string `env:"SERVER_PORT" env-default:"8080"`
	SigningKey  string `env:"SIGNING_KEY" env-required:"true"`
	DatabaseURL string `env:"DATABASE_URL" env-required:"true"`
	DocsHTML    string `env:"DOCS_HTML" env-required:"true"`
	DocsYAML    string `env:"DOCS_YAML" env-required:"true"`
}

func (c *Config) Addr() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func New() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg
}
