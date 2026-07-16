.DEFAULT_GOAL := help

# Colors
RED     := $(shell tput -Txterm setaf 1)
GREEN   := $(shell tput -Txterm setaf 2)
BLUE    := $(shell tput -Txterm setaf 4)
CYAN    := $(shell tput -Txterm setaf 6)
RESET   := $(shell tput -Txterm sgr0)

# Apps
SERVER_DIR=apps/server

# Docker
COMPOSE_PATH=compose.yaml

# Migration
MIGRATIONS_DIR := internal/database/mariadb/migrations
MIGRATE_CMD    := go run --tags mysql \
	github.com/golang-migrate/migrate/v4/cmd/migrate \
	-path=$(MIGRATIONS_DIR) \
	-database="mysql://$$DATABASE_USER:$$DATABASE_PASSWORD@tcp($$DATABASE_HOST:$$DATABASE_PORT)/$$DATABASE_NAME?multiStatements=true"

# Env
ENV := set -a; . "./.env"; set +a;

## Server 

.PHONY: server/run
server/run: ## Run the server application
	@$(ENV) cd $(SERVER_DIR) && go run cmd/main.go

## Migration

.PHONY: migrate/new
migrate/new: ## Create a new server migration. Usage: make migrate/new NAME=create_users
	@if [ -z "$(NAME)" ]; then \
		echo "ERROR: NAME is required. Example: make migrate/new NAME=create_users"; \
		exit 1; \
	fi
	@cd $(SERVER_DIR) && go run github.com/golang-migrate/migrate/v4/cmd/migrate \
		create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

.PHONY: migrate/up
migrate/up: ## Run server migrations
	@$(ENV) cd $(SERVER_DIR) && $(MIGRATE_CMD) up

.PHONY: migrate/down
migrate/down: ## Rollback one server migration
	@$(ENV) cd $(SERVER_DIR) && $(MIGRATE_CMD) down 1

.PHONY: migrate/drop
migrate/drop: ## Drop server database
	@$(ENV) cd $(SERVER_DIR) && $(MIGRATE_CMD) drop

.PHONY: migrate/version
migrate/version: ## Check server migration version
	@$(ENV) cd $(SERVER_DIR) && $(MIGRATE_CMD) version

## Go tooling

.PHONY: server/deps
deps: ## Download Go dependencies
	@cd $(SERVER_DIR) && go mod tidy

.PHONY: server/vuln
vuln: ## Check for vulnerabilities in Go dependencies
	@cd $(SERVER_DIR) && go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## Docker commands

.PHONY: docker/start
docker/start: ## Start all local containers.
	@$(ENV) docker compose -f ${COMPOSE_PATH} up -d

.PHONY: docker/stop
docker/stop: ## Stop all local containers.
	@docker compose -f ${COMPOSE_PATH} stop

.PHONY: docker/down
docker/down: ## Remove all local containers.
	@docker compose -f ${COMPOSE_PATH} down --remove-orphans

.PHONY: docker/reset
docker/reset: docker/down ## Stop, delete, build and start all local containers.
	@$(ENV) docker compose -f ${COMPOSE_PATH} up -d

## Help

.PHONY: help
help: ## Show this help
	@echo 'Usage:'
	@echo '  ${CYAN}make${RESET} ${GREEN}<target>${RESET}'
	@echo
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/[a-zA-Z_\-]+:.*?##.*$$/) {printf "    ${CYAN}%-25s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${BLUE}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
	@echo