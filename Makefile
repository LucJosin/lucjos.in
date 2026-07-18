.DEFAULT_GOAL := help
# Don't print the child makefiles while calling them
MAKEFLAGS += --no-print-directory

# Colors
RED     := $(shell tput -Txterm setaf 1)
GREEN   := $(shell tput -Txterm setaf 2)
YELLOW  := $(shell tput -Txterm setaf 3)
BLUE    := $(shell tput -Txterm setaf 4)
CYAN    := $(shell tput -Txterm setaf 6)
RESET   := $(shell tput -Txterm sgr0)

# Apps
SERVER_DIR=apps/server

# Docker
COMPOSE_PATH=compose.yaml

# Env
ENV := set -a; . "./.env"; set +a;

## Server 

.PHONY: server/run
server/run: ## Run the server application
	@$(MAKE) -C $(SERVER_DIR) server/run

## Quality and checks

.PHONY: lint
lint: ## Run linters
	@$(MAKE) -C $(SERVER_DIR) lint

.PHONY: lint/sql
lint/sql: ## Run SQL linter
	@$(MAKE) -C $(SERVER_DIR) lint/sql

.PHONY: fmt
fmt: ## Run Go and SQL formatters
	@$(MAKE) -C $(SERVER_DIR) fmt

.PHONY: test
test: ## Run tests
	@$(MAKE) -C $(SERVER_DIR) test

## Migration

.PHONY: migrate/new
migrate/new: ## Create a new server migration. Usage: make migrate/new NAME=create_users
	@$(MAKE) -C $(SERVER_DIR) migrate/new

.PHONY: migrate/up
migrate/up: ## Run server migrations
	$(MAKE) -C $(SERVER_DIR) migrate/up

.PHONY: migrate/down
migrate/down: ## Rollback one server migration
	@$(MAKE) -C $(SERVER_DIR) migrate/down

.PHONY: migrate/drop
migrate/drop: ## Drop server database
	@$(MAKE) -C $(SERVER_DIR) migrate/drop

.PHONY: migrate/version
migrate/version: ## Check server migration version
	@$(MAKE) -C $(SERVER_DIR) migrate/version

## Go tooling

.PHONY: deps
deps: ## Download Go dependencies
	@$(MAKE) -C $(SERVER_DIR) deps

.PHONY: vuln
vuln: ## Check for vulnerabilities in Go dependencies
	@$(MAKE) -C $(SERVER_DIR) vuln

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

.PHONY: validate-deps
validate-deps: ## Validate require dependencies
	@$(MAKE) -C $(SERVER_DIR) validate-deps

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
