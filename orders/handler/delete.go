package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sentrionic/ecommerce/common"
	"github.com/sentrionic/ecommerce/common/apperrors"
	status "github.com/sentrionic/ecommerce/common/order"
	"github.com/sentrionic/ecommerce/orders/ent"
	gen "github.com/sentrionic/ecommerce/orders/ent/order"
	"log"
	"net/http"
)

func (h *Handler) DeleteOrder(c *gin.Context) {
	userId := c.MustGet(common.UserID).(uuid.UUID)
	oid := c.Param("id")

	orderId, err := uuid.Parse(oid)

	if err != nil {
		e := apperrors.NewBadRequest("invalid id")
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
	}

	ctx := context.Background()

	order, err := h.db.Order.
		Query().
		WithProduct().
		Where(gen.IDEQ(orderId)).
		First(ctx)

	if err != nil {
		log.Println(err)
		if ent.IsNotFound(err) {
			e := apperrors.NewNotFound("order", oid)
			c.JSON(apperrors.Status(e), gin.H{
				"error": e,
			})
			return
		}
		e := apperrors.NewInternal()
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
		return
	}

	if order.UserID != userId {
		e := apperrors.NewAuthorization("you are not the owner")
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
		return
	}

	order.Status = status.Cancelled
	order.Version = order.Version + 1

	tx, err := h.db.Tx(ctx)
	if err != nil {
		log.Printf("failed creating transaction: %v", err)
		e := apperrors.NewInternal()
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
		return
	}

	if err = ent.UpdateOrderTx(tx, order); err != nil {
		log.Printf("unexpected failure: %v", err)
		err = tx.Rollback()
		log.Printf("error rolling back: %v", err)
		e := apperrors.NewInternal()
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
		return
	}

	if err != nil {
		log.Println(err)
		e := apperrors.NewInternal()
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
		return
	}

	h.p.PublishOrderCancelled(order)

	c.JSON(http.StatusOK, serializeOrderResponse(order))
}
