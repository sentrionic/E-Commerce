package publishers

import (
	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
	"github.com/sentrionic/ecommerce/common/events"
	"github.com/sentrionic/ecommerce/orders/ent"
)

type OrderPublisher interface {
	PublishOrderCreated(order *ent.Order, product *ent.Product)
	PublishOrderCancelled(order *ent.Order)
}

type orderPublisher struct {
	publisher *events.Publisher
}

func NewOrderPublisher(sc stan.Conn) OrderPublisher {
	return &orderPublisher{
		publisher: events.NewPublisher(sc),
	}
}

func (o *orderPublisher) PublishOrderCreated(order *ent.Order, product *ent.Product) {
	evt := events.OrderCreatedEvent(events.OrderCreatedData{
		ID:        order.ID,
		Status:    order.Status,
		UserId:    order.UserID,
		ExpiresAt: order.ExpiresAt,
		Product: struct {
			ID    uuid.UUID
			Price int
		}{
			ID:    product.ID,
			Price: product.Price,
		},
	})

	o.publisher.Publish(evt)
}

func (o *orderPublisher) PublishOrderCancelled(order *ent.Order) {
	evt := events.OrderCancelledEvent(events.OrderCancelledData{
		ID:        order.ID,
		ProductID: order.Edges.Product.ID,
	})

	o.publisher.Publish(evt)
}
