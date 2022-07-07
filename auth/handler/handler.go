package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sentrionic/ecommerce/auth/ent"
	"github.com/sentrionic/ecommerce/auth/utils"
	"github.com/sentrionic/ecommerce/common"
	"github.com/sentrionic/ecommerce/common/middleware"
	"github.com/sentrionic/ecommerce/common/token"
	"log"
	"time"
)

type Handler struct {
	db     *ent.Client
	config utils.Config
}

type Config struct {
	R      *gin.Engine
	DB     *ent.Client
	Config utils.Config
}

func NewHandler(c *Config) {
	h := &Handler{
		db:     c.DB,
		config: c.Config,
	}

	g := c.R.Group("/api/auth")

	g.GET("/current", middleware.AuthUser(h.config.SessionSecret), h.CurrentUser)
	g.POST("/register", h.Register)
	g.POST("/login", h.Login)
	g.POST("/logout", h.Logout)
}

// setUserSession saves the users ID in the session
func (h *Handler) setUserSession(c *gin.Context, id uuid.UUID) {
	accessToken, err := token.CreateToken(id, time.Hour, h.config.SessionSecret)

	if err != nil {
		log.Println("error creating access token")
		return
	}

	h.setCookie(c, accessToken, time.Hour.Seconds())
}

func (h *Handler) setCookie(c *gin.Context, accessToken string, maxAge float64) {
	c.SetCookie(
		common.Cookie,
		accessToken,
		int(maxAge),
		"/",
		h.config.Domain,
		gin.Mode() != gin.TestMode,
		true,
	)
}
