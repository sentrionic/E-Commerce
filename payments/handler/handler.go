package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/ecommerce/common/middleware"
	"github.com/sentrionic/ecommerce/payments/ent"
	"github.com/sentrionic/ecommerce/payments/publishers"
	"github.com/sentrionic/ecommerce/payments/service"
	"github.com/sentrionic/ecommerce/payments/utils"
)

type Handler struct {
	db     *ent.Client
	config utils.Config
	p      publishers.PaymentPublisher
	s      service.StripeService
}

type Config struct {
	R      *gin.Engine
	DB     *ent.Client
	Config utils.Config
	P      publishers.PaymentPublisher
	S      service.StripeService
}

func NewHandler(c *Config) {
	h := &Handler{
		db:     c.DB,
		config: c.Config,
		p:      c.P,
		s:      c.S,
	}

	g := c.R.Group("/api/payments")
	g.Use(middleware.AuthUser(h.config.SessionSecret))

	g.POST("", h.CreatePayment)
}
