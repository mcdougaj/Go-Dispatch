.PHONY: help setup dev-backend dev-frontend dev build test clean docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Initial setup - install dependencies
	@echo "Installing Go dependencies..."
	go mod download
	@echo "Installing frontend dependencies..."
	cd frontend && npm install
	@echo "Setup complete!"

dev-backend: ## Run backend in development mode
	go run cmd/api/main.go

dev-frontend: ## Run frontend in development mode
	cd frontend && npm start

dev: ## Run both backend and frontend (requires two terminals)
	@echo "Run 'make dev-backend' in one terminal and 'make dev-frontend' in another"

build: ## Build backend binary
	CGO_ENABLED=0 go build -o bin/go-dispatch ./cmd/api

build-frontend: ## Build frontend for production
	cd frontend && npm run build

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf frontend/build/
	rm -rf frontend/node_modules/

docker-up: ## Start all services with Docker Compose
	docker-compose up -d

docker-down: ## Stop all Docker services
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f

docker-rebuild: ## Rebuild and restart Docker containers
	docker-compose down
	docker-compose build
	docker-compose up -d

db-shell: ## Connect to PostgreSQL database
	docker-compose exec postgres psql -U dispatch_user -d motive_dispatch

redis-cli: ## Connect to Redis CLI
	docker-compose exec redis redis-cli
