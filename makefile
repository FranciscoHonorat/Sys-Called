MODULES := ./services/ticket-service ./services/employees-service
DOCKER_COMPOSE=docker compose

.PHONY: help build run test test-go test-frontend dev-frontend fmt tidy env docker-up docker-down docker-ps k8s-up k8s-deploy k8s-status k8s-down helm-lint up-all down-all logs

help:
	@echo "Available commands:"
	@echo "  build   - Build the Go modules"
	@echo "  vet     - Run go vet on the Go modules"
	@echo "  test    - Run the Go and frontend tests"
	@echo "  test-go - Run tests for the Go modules"
	@echo "  test-frontend - Run the frontend tests"
	@echo "  dev-frontend  - Start the frontend dev server (Vite)"
	@echo "  fmt     - Format the Go code in the services directories"
	@echo "  run     - Run the ticket-service"
	@echo "  env     - Create .env from .env.example with a fresh JWT signing key"
	@echo "  docker-up      - Build and start Docker containers in detached mode"
	@echo "  docker-down    - Stop Docker containers"
	@echo "  docker-ps      - Show the containers and their health"
	@echo "  k8s-up         - Create a kind cluster, build the images and install the Helm chart"
	@echo "  k8s-deploy     - Rebuild the images and upgrade the Helm release"
	@echo "  k8s-status     - Show pods, services and volumes in the cluster"
	@echo "  k8s-down       - Delete the kind cluster"
	@echo "  helm-lint      - Lint and render the Helm chart"
	@echo "  up-all         - Start both the Docker Compose stack and the kind cluster"
	@echo "  down-all       - Stop the Docker Compose stack and delete the kind cluster"
	@echo "  logs    - View logs for the ticket-service container"

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/ticket-service ./services/ticket-service/cmd/server
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/employees-service ./services/employees-service/cmd/server

run:
	go run ./$(SERVICE)/cmd/server

test: test-go test-frontend

test-go:
	@for m in $(MODULES); do \
		echo "==> go test ./... ($$m)"; \
		(cd $$m && go test ./...) || exit 1; \
	done

test-frontend:
	cd frontend && npm test

dev-frontend:
	cd frontend && npm run dev

fmt:
	@for m in $(MODULES); do \
		(cd $$m && go fmt ./...); \
	done

tidy:
	@for m in $(MODULES); do \
		echo "==> go mod tidy ($$m)"; \
		(cd $$m && go mod tidy) || exit 1; \
	done

env: .env

.env:
	@cp .env.example .env
	@{ printf 'JWT_PRIVATE_KEY="'; openssl genpkey -algorithm ed25519; printf '"\n'; } > .env.jwt
	@grep -v '^JWT_PRIVATE_KEY=' .env > .env.tmp && cat .env.tmp .env.jwt > .env && rm -f .env.tmp .env.jwt
	@echo "created .env"

docker-up: .env
	$(DOCKER_COMPOSE) up --build -d

docker-ps:
	$(DOCKER_COMPOSE) ps

docker-down:
	$(DOCKER_COMPOSE) down

k8s-up:
	./infra/scripts/k8s.sh up

k8s-deploy:
	./infra/scripts/k8s.sh deploy

k8s-status:
	./infra/scripts/k8s.sh status

k8s-down:
	./infra/scripts/k8s.sh down

helm-lint:
	helm lint infra/helm/sys-called
	helm template sys-called infra/helm/sys-called > /dev/null

up-all: docker-up k8s-up

down-all: docker-down k8s-down

logs:
	$(DOCKER_COMPOSE) logs -f ticket-service
