package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	gen "github.com/sentrionic/ecommerce/auth/ent/user"
	"github.com/sentrionic/ecommerce/common"
	"github.com/sentrionic/ecommerce/common/apperrors"
	"log"
	"net/http"
)

func (h *Handler) CurrentUser(c *gin.Context) {
	userId := c.MustGet(common.UserID).(uuid.UUID)

	user, err := h.db.User.Query().Where(gen.IDEQ(userId)).First(context.Background())

	if err != nil {
		log.Printf("Unable to find auth: %v\n%v", userId, err)
		e := apperrors.NewNotFound("userId", fmt.Sprintf("%d", userId))

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, serializeAuthResponse(user))
}
