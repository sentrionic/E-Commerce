package listeners

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
	"github.com/sentrionic/ecommerce/common/events"
	"github.com/sentrionic/ecommerce/orders/ent"
	gen "github.com/sentrionic/ecommerce/orders/ent/product"
	"log"
)

const queueGroupName = "orders-service"

type ProductListener interface {
	ProductCreatedListener(ctx context.Context)
	ProductUpdatedListener(ctx context.Context)
}

type productListener struct {
	listener *events.Listener
	client   *ent.Client
}

func NewProductListener(sc stan.Conn, client *ent.Client) ProductListener {
	return &productListener{
		listener: events.NewListener(sc, queueGroupName),
		client:   client,
	}
}

type productData struct {
	ID      uuid.UUID
	Version uint
	Title   string
	Price   int
}

func (p *productListener) ProductCreatedListener(ctx context.Context) {
	p.listener.Listen(events.ProductCreated, func(msg *stan.Msg) {
		var pd productData

		if err := json.Unmarshal(msg.Data, &pd); err != nil {
			log.Printf("error unmarshalling: %v", err)
			return
		}

		_, err := p.client.Product.Create().
			SetTitle(pd.Title).
			SetPrice(pd.Price).
			SetVersion(pd.Version).
			SetID(pd.ID).
			Save(ctx)

		if err != nil {
			log.Printf("error adding the product: %v", err)
			return
		}

		if err = msg.Ack(); err != nil {
			log.Printf("error acknowleding: %v", err)
		}
	})
}

func (p *productListener) ProductUpdatedListener(ctx context.Context) {
	p.listener.Listen(events.ProductUpdated, func(msg *stan.Msg) {
		var pd productData

		if err := json.Unmarshal(msg.Data, &pd); err != nil {
			log.Printf("error unmarshalling: %v", err)
		}

		product, err := p.client.Product.Query().
			Where(
				gen.And(
					gen.IDEQ(pd.ID),
					gen.VersionEQ(pd.Version-1),
				),
			).First(ctx)

		if err != nil {
			log.Printf("error getting the product: %v", err)
		}

		product.Title = pd.Title
		product.Price = pd.Price
		product.Version = pd.Version

		tx, err := p.client.Tx(ctx)
		if err != nil {
			log.Printf("failed creating transaction: %v", err)
			return
		}

		if err = ent.UpdateProductTx(tx, product); err != nil {
			log.Printf("unexpected failure: %v", err)
			return
		}

		if err = msg.Ack(); err != nil {
			log.Printf("error acknowleding: %v", err)
		}
	})
}
