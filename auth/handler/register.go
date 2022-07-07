package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/sentrionic/ecommerce/auth/ent"
	"github.com/sentrionic/ecommerce/common/apperrors"
	"github.com/sentrionic/ecommerce/common/middleware"
	"github.com/sentrionic/ecommerce/common/serializer"
	"log"
	"net/http"
	"strings"
)

type registerReq struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r registerReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.EmailFormat),
		validation.Field(&r.Username, validation.Required, validation.Length(3, 30)),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 150)),
	)
}

func (r *registerReq) sanitize() {
	r.Username = strings.TrimSpace(r.Username)
	r.Email = strings.TrimSpace(r.Email)
	r.Email = strings.ToLower(r.Email)
	r.Password = strings.TrimSpace(r.Password)
}

func (h *Handler) Register(c *gin.Context) {
	var req registerReq

	// Bind incoming json to struct and check for validation errors
	if ok := middleware.BindData(c, &req); !ok {
		return
	}

	req.sanitize()

	ctx := context.Background()
	user, err := h.db.User.
		Create().
		SetUsername(req.Username).
		SetEmail(req.Email).
		SetPassword(req.Password).
		Save(ctx)

	if err != nil {
		log.Println(err)
		if ent.IsConstraintError(err) {
			serializer.FieldError(c, "Email", apperrors.DuplicateEmail)
			return
		}

		err := apperrors.NewInternal()
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	h.setUserSession(c, user.ID)

	c.JSON(http.StatusCreated, serializeAuthResponse(user))
}
