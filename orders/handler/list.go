package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sentrionic/ecommerce/common"
	"github.com/sentrionic/ecommerce/common/apperrors"
	gen "github.com/sentrionic/ecommerce/orders/ent/order"
	"net/http"
)

func (h *Handler) GetOrders(c *gin.Context) {
	userId := c.MustGet(common.UserID).(uuid.UUID)

	orders, err := h.db.Order.
		Query().
		Where(gen.UserID(userId)).
		All(context.Background())

	if err != nil {
		e := apperrors.NewInternal()
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
		return
	}

	rsp := make([]OrderResponse, 0)

	for _, order := range orders {
		rsp = append(rsp, serializeOrderResponse(order))
	}

	c.JSON(http.StatusOK, rsp)
}
