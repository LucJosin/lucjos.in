.DEFAULT_GOAL := help

# Colors
RED     := $(shell tput -Txterm setaf 1)
GREEN   := $(shell tput -Txterm setaf 2)
YELLOW  := $(shell tput -Txterm setaf 3)
BLUE    := $(shell tput -Txterm setaf 4)
CYAN    := $(shell tput -Txterm setaf 6)
RESET   := $(shell tput -Txterm sgr0)

# Migration
MIGRATIONS_DIR := migrations
MIGRATE_CMD    := go run --tags mysql \
	github.com/golang-migrate/migrate/v4/cmd/migrate \
	-path=$(MIGRATIONS_DIR) \
	-database="mysql://$$DATABASE_USER:$$DATABASE_PASSWORD@tcp($$DATABASE_HOST:$$DATABASE_PORT)/$$DATABASE_NAME?multiStatements=true"

# Env
ENV := set -a; . "../../.env"; set +a;
GOLANGCI_LINT_VERSION := $(shell cat ../../.golangci-version)

## Server 

.PHONY: server/run
server/run: ## Run the server application
	@$(ENV) go run cmd/main.go

.PHONY: lint
lint: validate-deps lint/sql ## Run linters
	@go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_LINT_VERSION} run ./...

.PHONY: lint/sql
lint/sql: validate-deps ## Run SQL linter
	@sqlfluff lint $(MIGRATIONS_DIR)/*.sql

.PHONY: fmt
fmt: validate-deps ## Run Go and SQL formatters
	@go fmt ./... && sqlfluff fix $(MIGRATIONS_DIR)/*.sql --quiet

.PHONY: test
test: ## Run tests
	@go test ./... -timeout 5s -race -shuffle on

## Migration

.PHONY: migrate/new
migrate/new: ## Create a new server migration. Usage: make migrate/new NAME=create_users
	@if [ -z "$(NAME)" ]; then \
		echo "${RED}ERROR: NAME is required. Example: make migrate/new NAME=create_users"; \
		exit 1; \
	fi
	@go run github.com/golang-migrate/migrate/v4/cmd/migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

.PHONY: migrate/up
migrate/up: ## Run server migrations
	@$(ENV) && $(MIGRATE_CMD) up

.PHONY: migrate/down
migrate/down: ## Rollback one server migration
	@$(ENV) && $(MIGRATE_CMD) down 1

.PHONY: migrate/drop
migrate/drop: ## Drop server database
	@$(ENV) && $(MIGRATE_CMD) drop

.PHONY: migrate/version
migrate/version: ## Check server migration version
	@$(ENV) && $(MIGRATE_CMD) version

## Go tooling

.PHONY: deps
deps: ## Download Go dependencies
	@cd go mod tidy

.PHONY: vuln
vuln: ## Check for vulnerabilities in Go dependencies
	@cd go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## Help

.PHONY: validate-deps
validate-deps: ## Validate require dependencies
	@FAILED=0; \
	command -v go >/dev/null 2>&1 || { \
		echo "$(RED) Error: 'go' is not installed.$(RESET)"; \
		echo "   $(YELLOW) Go to https://go.dev/doc/install$(RESET)"; \
		FAILED=1; \
	}; \
	command -v sqlfluff >/dev/null 2>&1 || { \
		echo "$(RED) Error: 'sqlfluff' is not installed.$(RESET)"; \
		echo "   $(YELLOW) Go to https://docs.sqlfluff.com/en/stable/gettingstarted.html#installing-sqlfluff$(RESET)"; \
		FAILED=1; \
	}; \
	if [ $$FAILED -eq 1 ]; then \
		echo "\n$(RED) Please install the missing dependencies listed above.$(RESET)"; \
		exit 1; \
	fi

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
