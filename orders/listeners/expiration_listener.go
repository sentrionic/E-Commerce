package listeners

import (
	"context"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"github.com/sentrionic/ecommerce/common/events"
	status "github.com/sentrionic/ecommerce/common/order"
	"github.com/sentrionic/ecommerce/orders/ent"
	"github.com/sentrionic/ecommerce/orders/publishers"
	"log"
)

type ExpirationListener interface {
	ExpirationCompleteListener(ctx context.Context)
}

type expirationListener struct {
	listener  *events.Listener
	client    *ent.Client
	publisher publishers.OrderPublisher
}

func NewExpirationListener(sc stan.Conn, client *ent.Client, publisher publishers.OrderPublisher) ExpirationListener {
	return &expirationListener{
		listener:  events.NewListener(sc, queueGroupName),
		client:    client,
		publisher: publisher,
	}
}

func (p *expirationListener) ExpirationCompleteListener(ctx context.Context) {
	p.listener.Listen(events.ExpirationComplete, func(msg *stan.Msg) {
		var ed events.ExpirationCompleteData

		if err := json.Unmarshal(msg.Data, &ed); err != nil {
			log.Printf("error unmarshalling: %v", err)
			return
		}

		order, err := p.client.Order.Get(ctx, ed.OrderID)

		if err != nil {
			log.Printf("error getting the order: %v", err)
			if err = msg.Ack(); err != nil {
				log.Printf("error acknowleding: %v", err)
			}
			return
		}

		if order.Status == status.Complete {
			log.Printf("order is already completed")
			return
		}

		order.Status = status.Cancelled
		order.Version = order.Version + 1

		tx, err := p.client.Tx(ctx)
		if err != nil {
			log.Printf("failed creating transaction: %v", err)
			return
		}

		if err = ent.UpdateOrderTx(tx, order); err != nil {
			log.Printf("unexpected failure: %v", err)
			return
		}

		p.publisher.PublishOrderCancelled(order)

		if err = msg.Ack(); err != nil {
			log.Printf("error acknowleding: %v", err)
		}
	})
}
