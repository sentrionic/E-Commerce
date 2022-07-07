package middleware

import (
	"github.com/sentrionic/ecommerce/common/apperrors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Request contains the validate function which validates the request with bindData
type Request interface {
	Validate() error
}

// FieldError is used to help extract validation errors
type FieldError struct {
	// The property containing the error
	Field string `json:"field"`
	// The specific error message
	Message string `json:"message"`
} //@name FieldError

// BindData is helper function, returns false if data is not bound
func BindData(c *gin.Context, req Request) bool {
	// Bind incoming json to struct and check for validation errors
	if err := c.ShouldBind(req); err != nil {
		log.Printf("error binding data: %v", err)
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return false
	}

	if err := req.Validate(); err != nil {
		errors := strings.Split(err.Error(), ";")
		fErrors := make([]FieldError, 0)

		for _, e := range errors {
			split := strings.Split(e, ":")
			er := FieldError{
				Field:   strings.TrimSpace(split[0]),
				Message: strings.TrimSpace(split[1]),
			}
			fErrors = append(fErrors, er)
		}

		log.Printf("field errors: %v", fErrors)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": fErrors,
		})
		return false
	}
	return true
}
