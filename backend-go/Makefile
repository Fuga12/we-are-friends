COMPOSE ?= docker-compose

.PHONY: up down run test migrate sqlc

up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

run:
	go run ./cmd/api

test:
	go test ./...

migrate:
	go run ./cmd/api

sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.29.0 generate
