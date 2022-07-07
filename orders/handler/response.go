package handler

import (
	"github.com/sentrionic/ecommerce/common/order"
	"github.com/sentrionic/ecommerce/orders/ent"
	"time"
)

type OrderResponse struct {
	ID        string       `json:"id"`
	Status    order.Status `json:"order_status"`
	ExpiresAt time.Time    `json:"expires_at"`
	UserId    string       `json:"user_id"`
}

func serializeOrderResponse(o *ent.Order) OrderResponse {
	return OrderResponse{
		ID:        o.ID.String(),
		Status:    o.Status,
		ExpiresAt: o.ExpiresAt,
		UserId:    o.UserID.String(),
	}
}
