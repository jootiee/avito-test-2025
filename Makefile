include .env
export

DATABASE_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable
DOCKER_DATABASE_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:5432/$(POSTGRES_DB)?sslmode=disable
MIGRATE := $(shell which migrate || echo "$(HOME)/go/bin/migrate")

.PHONY: build
build:
	go build -v -o main ./cmd/app/main.go

.PHONY: test
test-unit:
	go test -v ./internal/...

.PHONY: test-unit-coverage
test-unit-coverage:
	go test -coverprofile=coverage.out ./internal/...
	go tool cover -func=coverage.out

.PHONY: test-unit-coverage-html
test-unit-coverage-html:
	go test -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

.PHONY: test-load
test-load: loadtest-team loadtest-pr loadtest-all

.PHONY: test-load-team
test-load-tea:
	@./loadtest/test-team.sh

.PHONY: test-load-pr
test-load-pr:
	@./loadtest/test-pr.sh

.PHONY: test-load-all
test-load-all:
	@./loadtest/test-all.sh

.PHONY: migrate-up
migrate-up:
	$(MIGRATE) -path migrations -database "$(DATABASE_URL)" up

.PHONY: migrate-down
migrate-down:
	$(MIGRATE) -path migrations -database "$(DATABASE_URL)" down

.PHONY: migrate-down-all
migrate-down-all:
	$(MIGRATE) -path migrations -database "$(DATABASE_URL)" down -all

.PHONY: migrate-force
migrate-force:
	@read -p "Enter version to force: " version; \
	$(MIGRATE) -path migrations -database "$(DATABASE_URL)" force $$version

.PHONY: migrate-version
migrate-version:
	$(MIGRATE) -path migrations -database "$(DATABASE_URL)" version

.PHONY: migrate-create
migrate-create:
	@read -p "Enter migration name: " name; \
	$(MIGRATE) create -ext sql -dir migrations -seq $$name

.PHONY: migrate-docker-up
migrate-docker-up:
	docker compose exec app migrate -path /app/migrations -database "$(DOCKER_DATABASE_URL)" up

.PHONY: migrate-docker-down
migrate-docker-down:
	docker compose exec app migrate -path /app/migrations -database "$(DOCKER_DATABASE_URL)" down

.PHONY: db-shell
db-shell:
	docker compose exec postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

.PHONY: clean
clean:
	rm *.out *.html

.DEFAULT_GOAL := build