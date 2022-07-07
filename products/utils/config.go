package utils

import (
	"context"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	DatabaseUrl   string `env:"DATABASE_URL"`
	SessionSecret string `env:"SESSION_SECRET"`
	NatsClientID  string `env:"NATS_CLIENT_ID"`
	NatsURL       string `env:"NATS_URL"`
	NatsClusterID string `env:"NATS_CLUSTER_ID"`
}

func LoadConfig(ctx context.Context) (config Config, err error) {
	err = envconfig.Process(ctx, &config)

	if err != nil {
		return
	}

	return
}
