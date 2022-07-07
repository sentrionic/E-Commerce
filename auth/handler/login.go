package handler

import (
	"context"
	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/sentrionic/ecommerce/auth/ent"
	gen "github.com/sentrionic/ecommerce/auth/ent/user"
	"github.com/sentrionic/ecommerce/common/apperrors"
	"github.com/sentrionic/ecommerce/common/middleware"
	"log"
	"net/http"
	"strings"
)

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r loginReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.EmailFormat),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 150)),
	)
}

func (r *loginReq) sanitize() {
	r.Email = strings.TrimSpace(r.Email)
	r.Email = strings.ToLower(r.Email)
	r.Password = strings.TrimSpace(r.Password)
}

func (h *Handler) Login(c *gin.Context) {
	var req loginReq

	if ok := middleware.BindData(c, &req); !ok {
		return
	}

	req.sanitize()

	ctx := context.Background()
	user, err := h.db.User.Query().Where(gen.EmailEQ(req.Email)).First(ctx)

	if err != nil {
		log.Println(err)
		if ent.IsNotFound(err) {
			err := apperrors.NewNotFound("email", req.Email)
			c.JSON(apperrors.Status(err), gin.H{
				"error": err,
			})
			return
		}
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	match, err := argon2id.ComparePasswordAndHash(req.Password, user.Password)

	if err != nil {
		err := apperrors.NewNotFound("email", req.Email)
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	if !match {
		err := apperrors.NewNotFound("email", req.Email)
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	h.setUserSession(c, user.ID)

	c.JSON(http.StatusOK, serializeAuthResponse(user))
}
