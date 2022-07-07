package publishers

import (
	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
	"github.com/sentrionic/ecommerce/common/events"
)

type ExpirationPublisher interface {
	PublishExpirationComplete(orderID uuid.UUID)
}

type expirationPublisher struct {
	publisher *events.Publisher
}

func NewExpirationPublisher(sc stan.Conn) ExpirationPublisher {
	return &expirationPublisher{
		publisher: events.NewPublisher(sc),
	}
}

func (e *expirationPublisher) PublishExpirationComplete(orderID uuid.UUID) {
	evt := events.ExpirationCompleteEvent(events.ExpirationCompleteData{
		OrderID: orderID,
	})

	e.publisher.Publish(evt)
}
