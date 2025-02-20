# Variables
APP_NAME=runner
DOCKER_COMPOSE=docker-compose

# Run the application
run:
	go run cmd/runner.go

# Build the Go binary
build:
	go build -o $(APP_NAME) cmd/runner.go

# Run Docker Compose
docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down



# Restart Docker services
docker-restart:
	$(DOCKER_COMPOSE) down && $(DOCKER_COMPOSE) up -d

# Check logs
docker-build:
	$(DOCKER_COMPOSE) up -d --build $(target)


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
