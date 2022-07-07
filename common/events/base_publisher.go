package events

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
)

type Publisher struct {
	client stan.Conn
}

func NewPublisher(client stan.Conn) *Publisher {
	return &Publisher{
		client: client,
	}
}

func (p *Publisher) Publish(event Event) {
	bytes, err := json.Marshal(event.Data)

	if err != nil {
		log.Printf("Error marshalling data: %v", err)
	}

	if err = p.client.Publish(string(event.Subject), bytes); err != nil {
		log.Printf("Error publishing event: %v", err)
	}

	log.Printf("Event published to subject %s", event.Subject)
}
