package utils

import (
	"entgo.io/ent/dialect"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sentrionic/ecommerce/auth/ent"
)

func SetupDatabase(config Config) (*ent.Client, error) {
	if gin.Mode() == gin.ReleaseMode {
		return ent.Open(dialect.Postgres, config.DatabaseUrl)
	}

	return ent.Open(dialect.Postgres, config.DatabaseUrl, ent.Debug())
}
