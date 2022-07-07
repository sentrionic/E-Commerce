package main

import (
	"context"
	"entgo.io/ent/dialect"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"github.com/sentrionic/ecommerce/payments/ent"
	_ "github.com/sentrionic/ecommerce/payments/ent/runtime"
	"github.com/sentrionic/ecommerce/payments/handler"
	"github.com/sentrionic/ecommerce/payments/listeners"
	"github.com/sentrionic/ecommerce/payments/publishers"
	"github.com/sentrionic/ecommerce/payments/service"
	"github.com/sentrionic/ecommerce/payments/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Starting server...")
	ctx := context.Background()

	config, err := utils.LoadConfig(ctx)

	if err != nil {
		log.Fatalln("Could not load the config")
	}

	client, err := ent.Open(dialect.Postgres, config.DatabaseUrl, ent.Debug())

	if err != nil {
		log.Fatalf("failed connecting to postgres: %v", err)
	}

	defer func(client *ent.Client) {
		err = client.Close()
		if err != nil {
			log.Fatalf("failed to close the db client: %v", err)
		}
	}(client)

	// Run migration.
	if err = client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	sc, err := stan.Connect(config.NatsClusterID, config.NatsClientID, stan.NatsURL(config.NatsURL))

	defer func(Client stan.Conn) {
		err = Client.Close()
	}(sc)

	if err != nil {
		log.Printf("error connecting to the nats client: %v", err)
	}

	publisher := publishers.NewPaymentPublisher(sc)

	ss := service.NewStripeService(config.StripeKey)

	router := gin.Default()

	handler.NewHandler(&handler.Config{
		R:      router,
		DB:     client,
		P:      publisher,
		Config: config,
		S:      ss,
	})

	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	listener := listeners.NewOrderListener(sc, client)
	listener.OrderCreatedListener(ctx)
	listener.OrderCancelledListener(ctx)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to initialize server: %v\n", err)
		}
	}()

	log.Printf("Listening on port %v\n", srv.Addr)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
}
