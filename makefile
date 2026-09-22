MODULES := ./services/ticket-service/... ./services/employees-service/... ./shared/...

.PHONY: build vet test fmt run up down logs

build:
	go build $(MODULES)

vet:
	go vet $(MODULES)

test:
	go test $(MODULES)

fmt:
	gofmt -l -w services shared

run:
	cd services/ticket-service && go run ./cmd/server

up:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f ticket-service
