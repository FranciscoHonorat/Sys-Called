MODULES := ./services/ticket-service ./services/employees-service ./shared
DOCKER_COMPOSE=docker compose

.PHONY: help build run test test-integration fmt tidy docker-build docker-build-inventory docker-build-cdc docker-build-user docker-up docker-down

help:
	@echo "Available commands:"
	@echo "  build   - Build the Go modules"
	@echo "  vet     - Run go vet on the Go modules"
	@echo "  test    - Run tests for the Go modules"
	@echo "  fmt     - Format the Go code in services and shared directories"
	@echo "  run     - Run the ticket-service"
	@echo "  docker-up      - Start Docker containers in detached mode"
	@echo "  docker-down    - Stop Docker containers"
	@echo "  logs    - View logs for the ticket-service container"

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/ticket-service ./services/ticket-service/cmd/server
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/employees-service ./services/employees-service/cmd/server

run:
	go run ./$(SERVICE)/cmd/server

test:
	@for m in $(MODULES); do \
		if [ -z "$$(find $$m -name '*.go' -print -quit)" ]; then \
			echo "==> skipping $$m (no Go files yet)"; \
			continue; \
		fi; \
		echo "==> go test ./... ($$m)"; \
		(cd $$m && go test ./...) || exit 1; \
	done

fmt:
	@for m in $(MODULES); do \
		(cd $$m && go fmt ./...); \
	done

tidy:
	@for m in $(MODULES); do \
		echo "==> go mod tidy ($$m)"; \
		(cd $$m && go mod tidy) || exit 1; \
	done

docker-up:
	$(DOCKER_COMPOSE) up --build -d

docker-down:
	$(DOCKER_COMPOSE) down

logs:
	$(DOCKER_COMPOSE) logs -f ticket-service
