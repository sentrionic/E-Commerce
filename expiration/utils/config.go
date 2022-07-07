package utils

import (
	"context"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	NatsClientID  string `env:"NATS_CLIENT_ID"`
	NatsURL       string `env:"NATS_URL"`
	NatsClusterID string `env:"NATS_CLUSTER_ID"`
	RedisHost     string `env:"REDIS_HOST"`
}

func LoadConfig(ctx context.Context) (config Config, err error) {
	err = envconfig.Process(ctx, &config)

	if err != nil {
		return
	}

	return
}
