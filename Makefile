include .env
export

export LOCAL_BIN:=$(CURDIR)/bin
export PATH:=$(LOCAL_BIN):$(PATH)

# HELP =================================================================================================================
# This will output the help for each task
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

run: ## Run the application
	go run cmd/app/main.go
.PHONY: run

migrate-up: ## Apply migrations
	migrate -path migrations -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" up
.PHONY: migrate-up

linter-golangci: ## Check Go code using golangci-lint
	golangci-lint run
.PHONY: linter-golangci

swag-v1: ## Initialize Swagger docs for v1
	swag init -g internal/controller/http/v1/router.go
.PHONY: swag-v1

up: ## Start Docker containers
	docker-compose up -d
.PHONY: up

