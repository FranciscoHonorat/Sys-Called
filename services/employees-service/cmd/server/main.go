package main

import (
	"context"
	"crypto/ed25519"
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
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/adapters/out/security"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application"
)

const (
	outboxRelayInterval = time.Second
	accessTokenTTL      = 15 * time.Minute
	refreshTokenTTL     = 7 * 24 * time.Hour
)

const seedPassword = "senha123"

var seedEmployees = []application.RegisterEmployeeInput{
	{ID: "agent-1", Name: "Ana Souza", Username: "ana", Password: seedPassword, Role: "support"},
	{ID: "agent-2", Name: "Bruno Lima", Username: "bruno", Password: seedPassword, Role: "support"},
	{ID: "agent-3", Name: "Carla Melo", Username: "carla", Password: seedPassword, Role: "support"},
	{ID: "admin-1", Name: "Administradora", Username: "admin", Password: seedPassword, Role: "admin"},
	{ID: "user-1", Name: "Usuário Padrão", Username: "usuario", Password: seedPassword, Role: "user"},
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

	signingKey, generated, err := security.LoadSigningKey(os.Getenv("JWT_PRIVATE_KEY"))
	if err != nil {
		log.Fatalf("failed to load JWT signing key: %v", err)
	}
	if generated {
		log.Println("JWT_PRIVATE_KEY not set: using an ephemeral signing key, tokens will not survive a restart")
	}

	repo := postgres.NewEmployeeRepository(pool)
	hasher := security.NewBcryptHasher()

	registerEmployee := application.NewRegisterEmployeeUseCase(repo, hasher)
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

	tokenIssuer := security.NewJWTIssuer(signingKey, accessTokenTTL)
	sessions := application.NewSessions(
		tokenIssuer,
		security.NewRefreshTokenGenerator(),
		postgres.NewRefreshTokenStore(pool),
		refreshTokenTTL,
		time.Now,
	)
	handler := httpapi.NewHandler(httpapi.UseCases{
		Authenticate:           application.NewAuthenticateUseCase(security.NewJWTVerifier(signingKey.Public().(ed25519.PublicKey))),
		ListEmployees:          application.NewListEmployeesUseCase(repo),
		Login:                  application.NewLoginUseCase(repo, hasher, sessions),
		RefreshSession:         application.NewRefreshSessionUseCase(repo, sessions),
		Logout:                 application.NewLogoutUseCase(sessions),
		PublicKeys:             application.NewGetPublicKeysUseCase(tokenIssuer),
		SignUp:                 application.NewSignUpUseCase(repo, hasher),
		RequestPasswordReset:   application.NewRequestPasswordResetUseCase(repo),
		ApproveEmployee:        application.NewApproveEmployeeUseCase(repo),
		IssueTemporaryPassword: application.NewIssueTemporaryPasswordUseCase(repo, hasher, security.NewTemporaryPasswordGenerator(), sessions),
		ChangePassword:         application.NewChangePasswordUseCase(repo, hasher),
	})
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
