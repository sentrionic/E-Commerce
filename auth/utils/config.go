package utils

import (
	"context"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	DatabaseUrl   string `env:"DATABASE_URL"`
	SessionSecret string `env:"SESSION_SECRET"`
	Domain        string `env:"DOMAIN"`
}

func LoadConfig(ctx context.Context) (config Config, err error) {
	err = envconfig.Process(ctx, &config)

	if err != nil {
		return
	}

	return
}
