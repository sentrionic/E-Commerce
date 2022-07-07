package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/ecommerce/auth/ent"
	_ "github.com/sentrionic/ecommerce/auth/ent/runtime"
	"github.com/sentrionic/ecommerce/auth/handler"
	"github.com/sentrionic/ecommerce/auth/utils"
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

	client, err := utils.SetupDatabase(config)

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

	router := gin.Default()

	handler.NewHandler(&handler.Config{
		R:  router,
		DB: client,
	})

	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

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
