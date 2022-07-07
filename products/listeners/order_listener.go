package listeners

import (
	"context"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"github.com/sentrionic/ecommerce/common/events"
	"github.com/sentrionic/ecommerce/products/ent"
	gen "github.com/sentrionic/ecommerce/products/ent/product"
	"github.com/sentrionic/ecommerce/products/publishers"
	"log"
)

const queueGroupName = "products-service"

type OrderListener interface {
	OrderCreatedListener(ctx context.Context)
	OrderCancelledListener(ctx context.Context)
}

type orderListener struct {
	listener  *events.Listener
	client    *ent.Client
	publisher publishers.ProductPublisher
}

func NewOrderListener(sc stan.Conn, client *ent.Client, publisher publishers.ProductPublisher) OrderListener {
	return &orderListener{
		listener:  events.NewListener(sc, queueGroupName),
		client:    client,
		publisher: publisher,
	}
}

func (o *orderListener) OrderCreatedListener(ctx context.Context) {
	o.listener.Listen(events.OrderCreated, func(msg *stan.Msg) {
		var od events.OrderCreatedData

		if err := json.Unmarshal(msg.Data, &od); err != nil {
			log.Printf("error unmarshalling: %v", err)
			return
		}

		product, err := o.client.Product.Query().Where(gen.IDEQ(od.Product.ID)).First(ctx)

		if err != nil {
			log.Printf("error getting the product: %v", err)
			return
		}

		product.OrderID = &od.ID
		product.Version = product.Version + 1

		tx, err := o.client.Tx(ctx)
		if err != nil {
			log.Printf("failed creating transaction: %v", err)
			return
		}
		if err = ent.UpdateProductTx(tx, product); err != nil {
			err = tx.Rollback()
			log.Printf("unexpected failure: %v", err)
			return
		}

		o.publisher.PublishProductUpdated(product)

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

		product, err := o.client.Product.Query().Where(gen.IDEQ(od.ProductID)).First(ctx)

		if err != nil {
			log.Printf("error getting the product: %v", err)
			return
		}

		product.OrderID = nil
		product.Version = product.Version + 1

		tx, err := o.client.Tx(ctx)
		if err != nil {
			log.Printf("failed creating transaction: %v", err)
			return
		}

		if err = ent.UpdateProductTx(tx, product); err != nil {
			err = tx.Rollback()
			log.Printf("unexpected failure: %v", err)
			return
		}

		o.publisher.PublishProductUpdated(product)

		if err = msg.Ack(); err != nil {
			log.Printf("error acknowleding: %v", err)
		}
	})
}
