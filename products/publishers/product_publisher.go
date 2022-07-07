package publishers

import (
	"github.com/nats-io/stan.go"
	"github.com/sentrionic/ecommerce/common/events"
	"github.com/sentrionic/ecommerce/products/ent"
)

type ProductPublisher interface {
	PublishProductCreated(product *ent.Product)
	PublishProductUpdated(product *ent.Product)
}

type productPublisher struct {
	publisher *events.Publisher
}

func NewProductPublisher(sc stan.Conn) ProductPublisher {
	return &productPublisher{
		publisher: events.NewPublisher(sc),
	}
}

func (p *productPublisher) PublishProductCreated(product *ent.Product) {
	evt := events.ProductCreatedEvent(events.ProductData{
		ID:      product.ID,
		Version: product.Version,
		Title:   product.Title,
		Price:   product.Price,
		UserId:  product.UserID,
	})

	p.publisher.Publish(evt)
}

func (p *productPublisher) PublishProductUpdated(product *ent.Product) {
	evt := events.ProductUpdatedEvent(events.ProductData{
		ID:      product.ID,
		Version: product.Version,
		Title:   product.Title,
		Price:   product.Price,
		UserId:  product.UserID,
	})

	p.publisher.Publish(evt)
}
