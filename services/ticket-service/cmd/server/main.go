package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/infra/cache"
	httpapi "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/infra/http"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/infra/messaging"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/infra/postgres"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		log.Fatal("KAFKA_BROKERS is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := postgres.NewPool(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	store := postgres.NewEventStore(pool)
	ticketCache := cache.NewInMemoryTicketCache()
	responsibles := postgres.NewResponsibleDirectory(pool)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: strings.Split(kafkaBrokers, ","),
		Topic:   "employees.events",
		GroupID: "ticket-service",
	})
	defer reader.Close()

	consumer := messaging.NewEmployeesConsumer(reader, responsibles)
	go consumer.Run(ctx)

	handler := httpapi.NewHandler(store, ticketCache, responsibles)
	router := httpapi.NewRouter(handler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("ticket-service listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}

	log.Println("server stopped")
}
