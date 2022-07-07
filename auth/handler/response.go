package handler

import "github.com/sentrionic/ecommerce/auth/ent"

type AuthResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func serializeAuthResponse(user *ent.User) *AuthResponse {
	return &AuthResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		Username: user.Username,
	}
}
