package publishers

import (
	"github.com/nats-io/stan.go"
	"github.com/sentrionic/ecommerce/common/events"
	"github.com/sentrionic/ecommerce/payments/ent"
)

type PaymentPublisher interface {
	PublishPaymentCreated(payment *ent.Payment)
}

type paymentPublisher struct {
	publisher *events.Publisher
}

func NewPaymentPublisher(sc stan.Conn) PaymentPublisher {
	return &paymentPublisher{
		publisher: events.NewPublisher(sc),
	}
}

func (o *paymentPublisher) PublishPaymentCreated(payment *ent.Payment) {
	evt := events.PaymentCreatedEvent(events.PaymentCreatedData{
		ID:       payment.ID,
		OrderID:  payment.OrderID,
		StripeID: payment.StripeID,
	})

	o.publisher.Publish(evt)
}
