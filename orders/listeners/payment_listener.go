package listeners

import (
	"context"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"github.com/sentrionic/ecommerce/common/events"
	status "github.com/sentrionic/ecommerce/common/order"
	"github.com/sentrionic/ecommerce/orders/ent"
	"log"
)

type PaymentListener interface {
	PaymentCreatedListener(ctx context.Context)
}

type paymentListener struct {
	listener *events.Listener
	client   *ent.Client
}

func NewPaymentListener(sc stan.Conn, client *ent.Client) PaymentListener {
	return &paymentListener{
		listener: events.NewListener(sc, queueGroupName),
		client:   client,
	}
}

func (p *paymentListener) PaymentCreatedListener(ctx context.Context) {
	p.listener.Listen(events.PaymentCreated, func(msg *stan.Msg) {
		var pd events.PaymentCreatedData

		if err := json.Unmarshal(msg.Data, &pd); err != nil {
			log.Printf("error unmarshalling: %v", err)
			return
		}

		order, err := p.client.Order.Get(ctx, pd.OrderID)

		if err != nil {
			log.Printf("error getting the order: %v", err)
			return
		}

		order.Status = status.Complete
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

		if err = msg.Ack(); err != nil {
			log.Printf("error acknowleding: %v", err)
		}
	})
}
