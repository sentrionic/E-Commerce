package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/sentrionic/ecommerce/common"
	"github.com/sentrionic/ecommerce/common/apperrors"
	"github.com/sentrionic/ecommerce/common/middleware"
	"github.com/sentrionic/ecommerce/common/order"
	"github.com/sentrionic/ecommerce/payments/ent"
	"log"
	"net/http"
)

type createRequest struct {
	OrderId string `json:"orderId"`
	Token   string `json:"token"`
}

func (r createRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.OrderId, validation.Required),
		validation.Field(&r.Token, validation.Required),
	)
}

func (h *Handler) CreatePayment(c *gin.Context) {
	userId := c.MustGet(common.UserID).(uuid.UUID)

	var req createRequest

	// Bind incoming json to struct and check for validation errors
	if ok := middleware.BindData(c, &req); !ok {
		return
	}

	orderId, err := uuid.Parse(req.OrderId)

	if err != nil {
		e := apperrors.NewBadRequest("invalid id")
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
	}

	ctx := context.Background()

	ord, err := h.db.Order.Get(ctx, orderId)

	if err != nil {
		log.Println(err)
		if ent.IsNotFound(err) {
			e := apperrors.NewNotFound("order", orderId.String())
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

	if ord.UserID != userId {
		e := apperrors.NewAuthorization("you are not the owner")
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
		return
	}

	if ord.Status == order.Cancelled {
		e := apperrors.NewBadRequest("cannot pay for a cancelled order")
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
		return
	}

	charge, err := h.s.HandleCharge(ord.Price, req.Token)

	if err != nil {
		log.Println(err)
		e := apperrors.NewInternal()
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
		return
	}

	payment, err := h.db.Payment.
		Create().
		SetOrderID(orderId).
		SetStripeID(charge.ID).
		Save(ctx)

	if err != nil {
		log.Println(err)
		e := apperrors.NewInternal()
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
		return
	}

	h.p.PublishPaymentCreated(payment)

	c.JSON(http.StatusCreated, serializePaymentResponse(payment))
}
