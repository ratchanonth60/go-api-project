# Variables
APP_NAME=runner
DOCKER_COMPOSE=docker-compose

# Run the application
run:
	go run cmd/runner.go

# Build the Go binary
build:
	go build -o $(APP_NAME) .

# Run Docker Compose
docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

# Restart Docker services
docker-restart:
	$(DOCKER_COMPOSE) down && $(DOCKER_COMPOSE) up -d

# Check logs
docker-logs:
	$(DOCKER_COMPOSE) logs -f

# Start specific services
docker-db:
	$(DOCKER_COMPOSE) up -d postgres

docker-redis:
	$(DOCKER_COMPOSE) up -d redis

# Stop specific services
docker-db-down:
	$(DOCKER_COMPOSE) stop postgres

docker-redis-down:
	$(DOCKER_COMPOSE) stop redis

# Run tests
test:
	go test ./...

# Format and lint code
fmt:
	go fmt ./...

lint:
	golangci-lint run# Clean build artifacts
clean:
	rm -f cmd/$(APP_NAME)

.PHONY: run build docker-up docker-down docker-restart docker-logs test fmt lint clea
