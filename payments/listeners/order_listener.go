package listeners

import (
	"context"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"github.com/sentrionic/ecommerce/common/events"
	status "github.com/sentrionic/ecommerce/common/order"
	"github.com/sentrionic/ecommerce/payments/ent"
	gen "github.com/sentrionic/ecommerce/payments/ent/order"
	"log"
)

const queueGroupName = "payments-service"

type OrderListener interface {
	OrderCreatedListener(ctx context.Context)
	OrderCancelledListener(ctx context.Context)
}

type orderListener struct {
	listener *events.Listener
	client   *ent.Client
}

func NewOrderListener(sc stan.Conn, client *ent.Client) OrderListener {
	return &orderListener{
		listener: events.NewListener(sc, queueGroupName),
		client:   client,
	}
}

func (o *orderListener) OrderCreatedListener(ctx context.Context) {
	o.listener.Listen(events.OrderCreated, func(msg *stan.Msg) {
		var od events.OrderCreatedData

		if err := json.Unmarshal(msg.Data, &od); err != nil {
			log.Printf("error unmarshalling: %v", err)
			return
		}

		_, err := o.client.Order.Create().
			SetID(od.ID).
			SetPrice(uint(od.Product.Price)).
			SetStatus(od.Status).
			SetUserID(od.UserId).
			SetVersion(od.Version).
			Save(ctx)

		if err != nil {
			log.Printf("error getting the product: %v", err)
			return
		}

		if err = msg.Ack(); err != nil {
			log.Printf("error acknowleding: %v", err)
		}
	})
}

func (o *orderListener) OrderCancelledListener(ctx context.Context) {
	o.listener.Listen(events.OrderCancelled, func(msg *stan.Msg) {
		var od events.OrderCancelledData

		if err := json.Unmarshal(msg.Data, &od); err != nil {
			log.Printf("error unmarshalling: %v", err)
			return
		}

		order, err := o.client.Order.Query().
			Where(
				gen.And(
					gen.IDEQ(od.ID),
					gen.VersionEQ(od.Version-1),
				),
			).First(ctx)

		if err != nil {
			log.Printf("error getting the product: %v", err)
			return
		}

		order.Status = status.Cancelled
		order.Version = od.Version

		tx, err := o.client.Tx(ctx)
		if err != nil {
			log.Printf("failed creating transaction: %v", err)
			return
		}

		if err = ent.UpdateOrderTx(tx, order); err != nil {
			err = tx.Rollback()
			log.Printf("unexpected failure: %v", err)
			return
		}

		if err = msg.Ack(); err != nil {
			log.Printf("error acknowleding: %v", err)
		}
	})
}
