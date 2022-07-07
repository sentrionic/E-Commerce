package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/ecommerce/common/middleware"
	"github.com/sentrionic/ecommerce/products/ent"
	"github.com/sentrionic/ecommerce/products/publishers"
	"github.com/sentrionic/ecommerce/products/utils"
)

type Handler struct {
	db     *ent.Client
	config utils.Config
	p      publishers.ProductPublisher
}

type Config struct {
	R      *gin.Engine
	DB     *ent.Client
	Config utils.Config
	P      publishers.ProductPublisher
}

func NewHandler(c *Config) {
	h := &Handler{
		db:     c.DB,
		config: c.Config,
		p:      c.P,
	}

	g := c.R.Group("/api/products")

	g.GET("", h.GetProducts)
	g.GET("/:id", h.GetProduct)

	g.Use(middleware.AuthUser(h.config.SessionSecret))

	g.POST("", h.CreateProduct)
	g.PUT("/:id", h.UpdateProduct)
}
