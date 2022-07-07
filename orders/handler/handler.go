package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/ecommerce/common/middleware"
	"github.com/sentrionic/ecommerce/orders/ent"
	"github.com/sentrionic/ecommerce/orders/publishers"
	"github.com/sentrionic/ecommerce/orders/utils"
)

type Handler struct {
	db     *ent.Client
	config utils.Config
	p      publishers.OrderPublisher
}

type Config struct {
	R      *gin.Engine
	DB     *ent.Client
	Config utils.Config
	P      publishers.OrderPublisher
}

func NewHandler(c *Config) {
	h := &Handler{
		db:     c.DB,
		config: c.Config,
		p:      c.P,
	}

	g := c.R.Group("/api/orders")
	g.Use(middleware.AuthUser(h.config.SessionSecret))

	g.GET("", h.GetOrders)
	g.POST("", h.CreateOrder)
	g.GET("/:id", h.GetOrder)
	g.DELETE("/:id", h.DeleteOrder)
}
