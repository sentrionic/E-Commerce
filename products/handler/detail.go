package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sentrionic/ecommerce/common/apperrors"
	"github.com/sentrionic/ecommerce/products/ent"
	gen "github.com/sentrionic/ecommerce/products/ent/product"
	"log"
	"net/http"
)

func (h *Handler) GetProduct(c *gin.Context) {
	pid := c.Param("id")

	productId, err := uuid.Parse(pid)

	if err != nil {
		e := apperrors.NewBadRequest("invalid id")
		c.JSON(apperrors.Status(e), gin.H{
			"error": e,
		})
	}

	product, err := h.db.Product.Query().Where(gen.IDEQ(productId)).First(context.Background())

	if err != nil {
		log.Println(err)
		if ent.IsNotFound(err) {
			e := apperrors.NewNotFound("product", pid)
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

	c.JSON(http.StatusOK, serializeProductResponse(product))
}
