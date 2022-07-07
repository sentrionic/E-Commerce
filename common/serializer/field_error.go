package serializer

import (
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/ecommerce/common/middleware"
	"net/http"
)

func FieldError(c *gin.Context, field, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"errors": []middleware.FieldError{
			{
				Field:   field,
				Message: message,
			},
		},
	})
}
