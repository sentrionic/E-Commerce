package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/ecommerce/common"
	"github.com/sentrionic/ecommerce/common/apperrors"
	"github.com/sentrionic/ecommerce/common/token"
)

// AuthUser checks if the request contains a valid session
// and saves the session's userId in the context
func AuthUser(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(common.Cookie)

		payload, err := token.VerifyToken(cookie, secret)
		if err != nil {
			e := apperrors.NewAuthorization(apperrors.InvalidSession)
			c.JSON(e.Status(), gin.H{
				"error": e,
			})
			c.Abort()
			return
		}

		id := payload.UID

		c.Set(common.UserID, id)

		c.Next()
	}
}
