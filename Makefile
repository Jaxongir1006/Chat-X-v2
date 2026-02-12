# =========================
# Project settings
# =========================
APP_NAME := chat-x-v2
MAIN := cmd/main.go

COMPOSE_FILE := compose-dev.yaml

# Goose settings
GOOSE_DRIVER := postgres
MIGRATIONS_DIR := migrations

# DB connection string
DATABASE_URL ?=

# Optional: load .env automatically if present (GNU Make include format)
# IMPORTANT: .env must be like KEY=value (no spaces)
ifneq (,$(wildcard .env))
	include .env
	export
endif

# =========================
# Helpers
# =========================
.PHONY: help debug
help:
	@echo ""
	@echo "Usage: make <target>"
	@echo ""
	@echo "App:"
	@echo "  make run             - Run HTTP server (go run cmd/main.go http)"
	@echo "  make superuser       - Create superuser (go run cmd/main.go superuser)"
	@echo ""
	@echo "Docker Compose (dev):"
	@echo "  make up              - docker compose up -d (compose-dev.yaml)"
	@echo "  make down            - docker compose down"
	@echo "  make restart         - down + up"
	@echo "  make logs            - docker compose logs -f"
	@echo "  make ps              - docker compose ps"
	@echo ""
	@echo "Goose migrations:"
	@echo "  make mig-status      - goose status"
	@echo "  make mig-up          - goose up"
	@echo "  make mig-down        - goose down"
	@echo "  make mig-reset       - goose reset"
	@echo "  make mig-create name=<migration_name> - create new migration"
	@echo ""
	@echo "Dev tools:"
	@echo "  make fmt             - gofmt"
	@echo "  make tidy            - go mod tidy"
	@echo "  make test            - go test ./..."
	@echo "  make vet             - go vet ./..."
	@echo ""

debug:
	@echo "DATABASE_URL=$(DATABASE_URL)"
	@echo "GOOSE_DRIVER=$(GOOSE_DRIVER)"
	@echo "MIGRATIONS_DIR=$(MIGRATIONS_DIR)"

# =========================
# App
# =========================
.PHONY: run superuser
run:
	go run $(MAIN) http

superuser:
	go run $(MAIN) superuser

# =========================
# Docker Compose (dev)
# =========================
.PHONY: up down restart logs ps
up:
	docker compose -f $(COMPOSE_FILE) up -d --build

down:
	docker compose -f $(COMPOSE_FILE) down

restart: down up

logs:
	docker compose -f $(COMPOSE_FILE) logs -f --tail=200

ps:
	docker compose -f $(COMPOSE_FILE) ps

# =========================
# Goose migrations (FIXED)
# Use env-mode for maximum compatibility:
# GOOSE_DRIVER=... GOOSE_DBSTRING=... goose COMMAND
# =========================
.PHONY: check-dsn mig-status mig-up mig-down mig-reset mig-create

check-dsn:
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "ERROR: DATABASE_URL is empty."; \
		echo "Put it in .env (project root) like:"; \
		echo "  DATABASE_URL=postgres://postgres:postgres@localhost:5432/chat_x?sslmode=disable"; \
		echo "Or export it:"; \
		echo "  export DATABASE_URL='postgres://postgres:postgres@localhost:5432/chat_x?sslmode=disable'"; \
		exit 1; \
	fi

mig-up: check-dsn
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(DATABASE_URL)" goose -dir $(MIGRATIONS_DIR) up

mig-down: check-dsn
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(DATABASE_URL)" goose -dir $(MIGRATIONS_DIR) down

mig-status: check-dsn
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(DATABASE_URL)" goose -dir $(MIGRATIONS_DIR) status

mig-reset: check-dsn
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(DATABASE_URL)" goose -dir $(MIGRATIONS_DIR) reset

mig-create:
	@if [ -z "$(name)" ]; then \
		echo "ERROR: migration name is required. Example:"; \
		echo "  make mig-create name=create_users_table"; \
		exit 1; \
	fi
	goose -dir $(MIGRATIONS_DIR) create "$(name)" sql

# =========================
# Dev tools
# =========================
.PHONY: fmt tidy test vet
fmt:
	gofmt -w .

tidy:
	go mod tidy

test:
	go test ./...

vet:
	go vet ./...

#-----------------------------------------#
###         Linting, formatting 		###
#-----------------------------------------#

.PHONY: lint-install
lint-install:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.5

.PHONY: lint
lint:
	golangci-lint run --max-issues-per-linter=0 --max-same-issues=0 ./...
