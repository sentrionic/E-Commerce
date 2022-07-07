package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/ecommerce/common"

	"net/http"
)

func (h *Handler) Logout(c *gin.Context) {
	c.Set(common.UserID, nil)

	h.setCookie(c, "", -1)

	c.JSON(http.StatusOK, true)
}
