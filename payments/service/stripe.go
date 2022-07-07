package service

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
)

type StripeService interface {
	HandleCharge(amount uint, token string) (*stripe.Charge, error)
}

type stripeService struct {
	sc client.API
}

func NewStripeService(key string) StripeService {
	sc := client.API{}
	sc.Init(key, nil)
	return stripeService{
		sc: sc,
	}
}

func (s stripeService) HandleCharge(amount uint, token string) (*stripe.Charge, error) {
	return s.sc.Charges.New(&stripe.ChargeParams{
		Amount:   stripe.Int64(int64(amount * 100)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Source:   &stripe.SourceParams{Token: stripe.String(token)},
	})
}
