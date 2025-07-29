.PHONY: help build run dev stop clean logs test

# Default target
help:
	@echo "Available commands:"
	@echo "  build     - Build the production Docker image"
	@echo "  run       - Run the production application"
	@echo "  dev       - Run the development environment with hot reload"
	@echo "  stop      - Stop all running containers"
	@echo "  clean     - Remove all containers and images"
	@echo "  logs      - Show logs from running containers"
	@echo "  test      - Run tests"
	@echo "  db-shell  - Connect to CockroachDB shell"
	@echo "  frontend  - Run only frontend in development mode"
	@echo "  backend   - Run only backend in development mode"

# Production build and run
build:
	@echo "Building production Docker image..."
	docker compose build

run:
	@echo "Starting production application..."
	docker compose up -d
	@echo "Application available at:"
	@echo "  Frontend: http://localhost"
	@echo "  Backend API: http://localhost:8000"
	@echo "  CockroachDB UI: http://localhost:8080"

# Development environment
dev:
	@echo "Starting development environment..."
	docker compose -f docker-compose.dev.yaml up --build

dev-detached:
	@echo "Starting development environment in background..."
	docker compose -f docker-compose.dev.yaml up -d --build
	@echo "Development environment available at:"
	@echo "  Frontend: http://localhost:5173"
	@echo "  Backend API: http://localhost:8000"
	@echo "  CockroachDB UI: http://localhost:8081"

# Individual services for development
frontend:
	@echo "Starting frontend development server..."
	docker compose -f docker-compose.dev.yaml up frontend --build

backend:
	@echo "Starting backend development server..."
	docker compose -f docker-compose.dev.yaml up backend crdb db-init --build

# Management commands
stop:
	@echo "Stopping all containers..."
	docker compose down
	docker compose -f docker-compose.dev.yaml down

clean:
	@echo "Cleaning up containers and images..."
	docker compose down -v --rmi all
	docker compose -f docker-compose.dev.yaml down -v --rmi all
	docker system prune -f

logs:
	@echo "Showing logs..."
	docker compose logs -f

logs-dev:
	@echo "Showing development logs..."
	docker compose -f docker-compose.dev.yaml logs -f

# Database operations
db-shell:
	@echo "Connecting to CockroachDB shell..."
	docker compose exec crdb cockroach sql --insecure

db-shell-dev:
	@echo "Connecting to development CockroachDB shell..."
	docker compose -f docker-compose.dev.yaml exec crdb cockroach sql --insecure

# Testing
test:
	@echo "Running tests..."
	@echo "Frontend tests:"
	docker compose -f docker-compose.dev.yaml exec frontend npm test
	@echo "Backend tests:"
	docker compose -f docker-compose.dev.yaml exec backend go test ./...

# Health check
health:
	@echo "Checking application health..."
	@curl -f http://localhost/health || echo "Application not responding"
	@curl -f http://localhost:8000/api/v1/health || echo "Backend API not responding"

# Quick setup for new developers
setup:
	@echo "Setting up development environment..."
	@echo "Building development images..."
	make dev-detached
	@echo "Waiting for services to start..."
	sleep 30
	@echo "Setup complete! Check the services:"
	make health