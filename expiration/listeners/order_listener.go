package listeners

import (
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/nats-io/stan.go"
	"github.com/sentrionic/ecommerce/common/events"
	"github.com/sentrionic/ecommerce/expiration/tasks"
	"log"
	"time"
)

const queueGroupName = "expiration-service"

type OrderListener interface {
	OrderCreatedListener()
}

type orderListener struct {
	listener  *events.Listener
	scheduler tasks.OrderScheduler
	client    *asynq.Client
}

func NewOrderListener(sc stan.Conn, client *asynq.Client, scheduler tasks.OrderScheduler) OrderListener {
	return &orderListener{
		listener:  events.NewListener(sc, queueGroupName),
		client:    client,
		scheduler: scheduler,
	}
}

func (o *orderListener) OrderCreatedListener() {
	o.listener.Listen(events.OrderCreated, func(msg *stan.Msg) {
		var od events.OrderCreatedData

		if err := json.Unmarshal(msg.Data, &od); err != nil {
			log.Printf("error unmarshalling: %v", err)
			return
		}

		delay := time.Until(od.ExpiresAt)

		fmt.Printf("Waiting this many milliseconds to process the job: %v", delay)

		task, err := o.scheduler.PublishExpirationTask(od.ID)

		if err != nil {
			log.Printf("error getting the task: %v", err)
			return
		}

		_, err = o.client.Enqueue(task, asynq.ProcessIn(delay))

		if err != nil {
			log.Printf("error getting the queue: %v", err)
			return
		}

		if err = msg.Ack(); err != nil {
			log.Printf("error acknowleding: %v", err)
		}
	})
}
