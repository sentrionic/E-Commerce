package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/sentrionic/ecommerce/common"
	"github.com/sentrionic/ecommerce/common/apperrors"
	"github.com/sentrionic/ecommerce/common/middleware"
	"log"
	"net/http"
	"strings"
)

type createRequest struct {
	Title string `json:"title"`
	Price int    `json:"price"`
}

func (r createRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Title, validation.Required),
		validation.Field(&r.Price, validation.Required, validation.Min(1)),
	)
}

func (r *createRequest) sanitize() {
	r.Title = strings.TrimSpace(r.Title)
}

func (h *Handler) CreateProduct(c *gin.Context) {
	var req createRequest

	// Bind incoming json to struct and check for validation errors
	if ok := middleware.BindData(c, &req); !ok {
		return
	}

	req.sanitize()

	userId := c.MustGet(common.UserID).(uuid.UUID)

	ctx := context.Background()
	product, err := h.db.Product.
		Create().
		SetTitle(req.Title).
		SetPrice(req.Price).
		SetUserID(userId).
		Save(ctx)

	if err != nil {
		log.Println(err)
		e := apperrors.NewInternal()
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
		return
	}

	h.p.PublishProductCreated(product)

	c.JSON(http.StatusCreated, serializeProductResponse(product))
}
