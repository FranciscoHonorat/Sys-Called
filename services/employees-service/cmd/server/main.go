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

	httpapi "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/adapters/in/http"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/adapters/in/outbox"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/adapters/out/messaging"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/adapters/out/postgres"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application"
)

const outboxRelayInterval = time.Second

var seedEmployees = []application.RegisterEmployeeInput{
	{ID: "agent-1", Name: "Ana Souza"},
	{ID: "agent-2", Name: "Bruno Lima"},
	{ID: "agent-3", Name: "Carla Melo"},
}

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
		port = "8081"
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := postgres.NewPool(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	repo := postgres.NewEmployeeRepository(pool)

	registerEmployee := application.NewRegisterEmployeeUseCase(repo)
	for _, seed := range seedEmployees {
		if err := registerEmployee.Execute(ctx, seed); err != nil {
			log.Fatalf("failed to seed employee %s: %v", seed.ID, err)
		}
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(strings.Split(kafkaBrokers, ",")...),
		Topic:    "employees.events",
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	relay := outbox.NewRelay(application.NewPublishPendingEventsUseCase(postgres.NewOutboxStore(pool), messaging.NewKafkaPublisher(writer)))
	go relay.Run(ctx, outboxRelayInterval)

	handler := httpapi.NewHandler(application.NewListEmployeesUseCase(repo))
	router := httpapi.NewRouter(handler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("employees-service listening on :%s", port)
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
