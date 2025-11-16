.PHONY: build up down lint lint-fix test-unit test-unit-coverage test-unit-coverage-html test-load test-load-team test-load-pr test-load-all test migrate-up migrate-down migrate-down-all migrate-force migrate-version migrate-create migrate-docker-up migrate-docker-down db-shell clean
include .env
export

DATABASE_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable
DOCKER_DATABASE_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:5432/$(POSTGRES_DB)?sslmode=disable
MIGRATE := $(shell which migrate || echo "$(HOME)/go/bin/migrate")
BINARY := main

help:
	@echo "Available commands:"
	@echo "- help: show this message"
	@echo ""
	@echo "- build: Build the Go application"
	@echo "- up: Run app in Docker container and run migrations"
	@echo "- down: Stop app in Docker container"
	@echo ""
	@echo "- lint: Run golangci-lint on the codebase"
	@echo "- lint-fix: Run golangci-lint with auto-fix on the codebase"
	@echo ""
	@echo "- test-unit: Run unit tests"
	@echo "- test-unit-coverage: Run unit tests with coverage report"
	@echo "- test-unit-coverage-html: Generate HTML coverage report"
	@echo "- test-load: Run all load tests"
	@echo "- test-load-team: Run team load test"
	@echo "- test-load-pr: Run PR load test"
	@echo "- test-load-all: Run all load tests"
	@echo "- test: Run all tests (load and unit with coverage)"
	@echo ""
	@echo "- migrate-up: Apply all up migrations"
	@echo "- migrate-down: Apply one down migration"
	@echo "- migrate-down-all: Apply all down migrations"
	@echo "- migrate-force: Force set migration version"
	@echo "- migrate-version: Show current migration version"

build:
	go build -v -o $(BINARY) ./cmd/app/main.go

up:
	docker compose --env-file .env -f deploy/docker-compose.yml up -d --wait
	$(MAKE) migrate-docker-up

down:
	docker compose -f deploy/docker-compose.yml down

lint:
	$(shell go env GOPATH)/bin/golangci-lint run ./...

lint-fix:
	$(shell go env GOPATH)/bin/golangci-lint run --fix ./...


test-unit:
	go test -v -race ./internal/service/... ./internal/handler/...

test-unit-coverage:
	go test -race -coverprofile=coverage.out ./internal/service/... ./internal/handler/...
	go tool cover -func=coverage.out

test-unit-coverage-html:
	go test -race -coverprofile=coverage.out ./internal/service/... ./internal/handler/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

test-load: test-load-team test-load-pr test-load-all

test-load-team:
	@./loadtest/test-team.sh

test-load-pr:
	@./loadtest/test-pr.sh

test-load-all:
	@./loadtest/test-all.sh

test: test-load-all test-unit-coverage


migrate-up:
	$(MIGRATE) -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	$(MIGRATE) -path migrations -database "$(DATABASE_URL)" down

migrate-down-all:
	$(MIGRATE) -path migrations -database "$(DATABASE_URL)" down -all

migrate-force:
	@read -p "Enter version to force: " version; \
	$(MIGRATE) -path migrations -database "$(DATABASE_URL)" force $$version

migrate-version:
	$(MIGRATE) -path migrations -database "$(DATABASE_URL)" version

migrate-create:
	@read -p "Enter migration name: " name; \
	$(MIGRATE) create -ext sql -dir migrations -seq $$name

migrate-docker-up:
	docker compose -f deploy/docker-compose.yml exec app migrate -path /app/migrations -database "$(DOCKER_DATABASE_URL)" up

migrate-docker-down:
	docker compose -f deploy/docker-compose.yml exec app migrate -path /app/migrations -database "$(DOCKER_DATABASE_URL)" down

db-shell:
	docker compose -f deploy/docker-compose.yml exec postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

clean:
	rm *.out *.html $(BINARY)

.DEFAULT_GOAL := help