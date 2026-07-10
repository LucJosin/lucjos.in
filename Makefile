.DEFAULT_GOAL := help

# Colors
RED     := $(shell tput -Txterm setaf 1)
GREEN   := $(shell tput -Txterm setaf 2)
BLUE    := $(shell tput -Txterm setaf 4)
CYAN    := $(shell tput -Txterm setaf 6)
RESET   := $(shell tput -Txterm sgr0)

# Env
ENV := set -a; . "./.env"; set +a;

## Server

.PHONY: run
run: ## Run the server application
	@$(ENV) go run apps/server/cmd/main.go

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