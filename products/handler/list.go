package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/ecommerce/common/apperrors"
	"log"
	"net/http"
)

func (h *Handler) GetProducts(c *gin.Context) {
	products, err := h.db.Product.Query().All(context.Background())

	if err != nil {
		log.Println(err)
		e := apperrors.NewInternal()
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
		return
	}

	rsp := make([]ProductResponse, 0)

	for _, product := range products {
		rsp = append(rsp, serializeProductResponse(product))
	}

	c.JSON(http.StatusOK, rsp)
}
