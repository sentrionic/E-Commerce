package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sentrionic/ecommerce/common"
	"github.com/sentrionic/ecommerce/common/apperrors"
	"github.com/sentrionic/ecommerce/common/middleware"
	"github.com/sentrionic/ecommerce/common/order"
	"github.com/sentrionic/ecommerce/orders/ent"
	ogen "github.com/sentrionic/ecommerce/orders/ent/order"
	gen "github.com/sentrionic/ecommerce/orders/ent/product"
	"log"
	"net/http"
	"time"
)

type createRequest struct {
	ProductId string `json:"productId"`
}

func (r createRequest) Validate() error {
	return nil
}

func (h *Handler) CreateOrder(c *gin.Context) {
	userId := c.MustGet(common.UserID).(uuid.UUID)

	var req createRequest

	// Bind incoming json to struct and check for validation errors
	if ok := middleware.BindData(c, &req); !ok {
		return
	}

	productId, err := uuid.Parse(req.ProductId)

	if err != nil {
		e := apperrors.NewBadRequest("invalid id")
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
	}

	ctx := context.Background()

	product, err := h.db.Product.
		Query().
		Where(gen.IDEQ(productId)).
		First(ctx)

	if err != nil {
		log.Println(err)
		if ent.IsNotFound(err) {
			e := apperrors.NewNotFound("product", productId.String())
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

	isReserved, err := h.db.Order.
		Query().
		Where(
			ogen.And(
				ogen.HasProductWith(gen.IDEQ(productId)),
				ogen.StatusIn(order.Created, order.AwaitingPayment, order.Complete),
			),
		).
		Exist(ctx)

	if err != nil {
		log.Println(err)
		e := apperrors.NewInternal()
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
		return
	}

	if isReserved {
		e := apperrors.NewBadRequest("Product is already reserved")
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
		return
	}

	newOrder, err := h.db.Order.
		Create().
		SetProduct(product).
		SetUserID(userId).
		SetStatus(order.Created).
		SetExpiresAt(time.Now().Add(time.Minute * 15)).
		Save(ctx)

	if err != nil {
		log.Println(err)
		e := apperrors.NewInternal()
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
		return
	}

	h.p.PublishOrderCreated(newOrder, product)

	c.JSON(http.StatusCreated, serializeOrderResponse(newOrder))
}
