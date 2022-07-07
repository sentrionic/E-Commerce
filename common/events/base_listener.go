package events

import (
	"github.com/nats-io/stan.go"
	"log"
	"time"
)

type Listener struct {
	client    stan.Conn
	queueName string
}

type OnMessage func(msg *stan.Msg)

func NewListener(client stan.Conn, queueName string) *Listener {
	return &Listener{
		client:    client,
		queueName: queueName,
	}
}

func (p *Listener) Listen(subject Subject, cb OnMessage) {
	_, err := p.client.QueueSubscribe(string(subject), p.queueName, func(msg *stan.Msg) {
		log.Printf("Message received: %s / %s", subject, p.queueName)
		cb(msg)
	},
		stan.SetManualAckMode(),
		stan.DeliverAllAvailable(),
		stan.AckWait(time.Second*5),
		stan.DurableName(p.queueName),
	)

	if err != nil {
		log.Printf("Error listening : %v", err)
	}
}
