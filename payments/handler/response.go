package handler

import (
	"github.com/sentrionic/ecommerce/payments/ent"
)

type PaymentResponse struct {
	ID string `json:"id"`
}

func serializePaymentResponse(payment *ent.Payment) PaymentResponse {
	return PaymentResponse{
		ID: payment.ID.String(),
	}
}
