FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/server ./cmd/server

FROM alpine:3.24

RUN apk add --no-cache ca-certificates && \
    adduser -D -u 10001 appuser

COPY --from=builder /out/server /usr/local/bin/server
COPY --from=builder /src/internal/adapters/out/postgres/migrations /migrations

USER appuser

EXPOSE 8081

ENTRYPOINT ["/usr/local/bin/server"]
