package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sentrionic/ecommerce/common"
	"github.com/sentrionic/ecommerce/common/apperrors"
	"github.com/sentrionic/ecommerce/orders/ent"
	gen "github.com/sentrionic/ecommerce/orders/ent/order"
	"log"
	"net/http"
)

func (h *Handler) GetOrder(c *gin.Context) {
	oid := c.Param("id")
	userId := c.MustGet(common.UserID).(uuid.UUID)

	orderId, err := uuid.Parse(oid)

	if err != nil {
		e := apperrors.NewBadRequest("invalid id")
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
	}

	order, err := h.db.Order.Query().Where(gen.IDEQ(orderId)).First(context.Background())

	if err != nil {
		log.Println(err)
		if ent.IsNotFound(err) {
			e := apperrors.NewNotFound("order", oid)
			c.JSON(apperrors.Status(e), gin.H{
				"error": e,
			})
		}
		e := apperrors.NewInternal()
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
		return
	}

	if order.UserID != userId {
		e := apperrors.NewAuthorization("you are not the order owner")
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, serializeOrderResponse(order))
}
