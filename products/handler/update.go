package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/sentrionic/ecommerce/common"
	"github.com/sentrionic/ecommerce/common/apperrors"
	"github.com/sentrionic/ecommerce/common/middleware"
	"github.com/sentrionic/ecommerce/products/ent"
	gen "github.com/sentrionic/ecommerce/products/ent/product"
	"log"
	"net/http"
	"strings"
)

type updateRequest struct {
	Title string `json:"title"`
	Price int    `json:"price"`
}

func (r updateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Title, validation.Required),
		validation.Field(&r.Price, validation.Required, validation.Min(1)),
	)
}

func (r *updateRequest) sanitize() {
	r.Title = strings.TrimSpace(r.Title)
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	userId := c.MustGet(common.UserID).(uuid.UUID)
	pid := c.Param("id")

	productId, err := uuid.Parse(pid)

	if err != nil {
		e := apperrors.NewBadRequest("invalid id")
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
	}

	var req updateRequest

	// Bind incoming json to struct and check for validation errors
	if ok := middleware.BindData(c, &req); !ok {
		return
	}

	req.sanitize()

	ctx := context.Background()

	product, err := h.db.Product.
		Query().
		Where(gen.IDEQ(productId)).
		First(ctx)

	if err != nil {
		log.Println(err)
		if ent.IsNotFound(err) {
			e := apperrors.NewNotFound("product", pid)
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

	if product.OrderID != nil {
		e := apperrors.NewBadRequest("cannot edit a reserved ticket")
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
		return
	}

	if product.UserID != userId {
		e := apperrors.NewAuthorization("you are not the owner")
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
		return
	}

	product.Title = req.Title
	product.Price = req.Price
	product.Version = product.Version + 1

	tx, err := h.db.Tx(ctx)
	if err != nil {
		log.Printf("failed creating transaction: %v", err)
		e := apperrors.NewInternal()
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
		return
	}

	if err = ent.UpdateProductTx(tx, product); err != nil {
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

	h.p.PublishProductUpdated(product)

	c.JSON(http.StatusOK, serializeProductResponse(product))
}
