.PHONY: help build run test vet lint clean tidy redis-up redis-down docker-build docker-up docker-down

APP_NAME := rate-limiter
CMD_DIR := ./cmd/server

help: ## Show available commands
	@echo "Usage: make [target]\n"
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	go build -o $(APP_NAME) $(CMD_DIR)

run: build ## Build and run the server locally
	./$(APP_NAME)

test: ## Run all tests
	go test ./...

vet: ## Run go vet static analysis
	go vet ./...

lint: ## Run golangci-lint
	golangci-lint run ./...

clean: ## Remove binary and build cache
	rm -f $(APP_NAME)
	go clean

tidy: ## Tidy go module dependencies
	go mod tidy

redis-up: ## Start only Redis
	docker compose up -d redis

redis-down: ## Stop only Redis
	docker compose stop redis

docker-build: ## Build Docker image
	docker compose build

docker-up: ## Start Redis and app with Docker Compose
	docker compose up --build

docker-down: ## Stop and remove Docker containers
	docker compose down

.DEFAULT_GOAL := help
