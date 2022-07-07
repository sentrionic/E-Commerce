package main

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/nats-io/stan.go"
	"github.com/sentrionic/ecommerce/expiration/listeners"
	"github.com/sentrionic/ecommerce/expiration/publishers"
	"github.com/sentrionic/ecommerce/expiration/tasks"
	"github.com/sentrionic/ecommerce/expiration/utils"
	"log"
)

func main() {
	log.Println("Starting Expiration Server...")
	ctx := context.Background()

	config, err := utils.LoadConfig(ctx)

	if err != nil {
		log.Fatalln("Could not load the config")
	}

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: config.RedisHost},
		asynq.Config{},
	)

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: config.RedisHost})
	defer client.Close()

	if err != nil {
		log.Printf("error connecting to the redis queue: %v", err)
	}

	sc, err := stan.Connect(config.NatsClusterID, config.NatsClientID, stan.NatsURL(config.NatsURL))

	defer func(Client stan.Conn) {
		err = Client.Close()
	}(sc)

	if err != nil {
		log.Printf("error connecting to the nats client: %v", err)
	}

	publisher := publishers.NewExpirationPublisher(sc)
	scheduler := tasks.NewOrderListener(publisher)

	listener := listeners.NewOrderListener(sc, client, scheduler)
	listener.OrderCreatedListener()

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeOrderExpiration, scheduler.HandleExpirationTask)

	if err = srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
