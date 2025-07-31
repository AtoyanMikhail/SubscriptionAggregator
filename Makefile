.PHONY: build run test clean swagger docker-up docker-down migrate-up migrate-down

# Build the application
build:
	go build -o bin/subscription-service cmd/main.go

# Run the application
run:
	go run cmd/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf docs/

# Generate Swagger documentation
swagger:
	swag init -g cmd/main.go -o docs --parseDependency --parseInternal

# Install swag tool if not present
install-swag:
	go install github.com/swaggo/swag/cmd/swag@latest

# Start PostgreSQL with docker-compose
docker-up:
	docker-compose up -d postgres

# Stop docker containers
docker-down:
	docker-compose down

# Run database migrations up
migrate-up:
	migrate -path migrations -database "postgres://postgres:password@localhost:5432/jwt?sslmode=disable" up

# Run database migrations down
migrate-down:
	migrate -path migrations -database "postgres://postgres:password@localhost:5432/jwt?sslmode=disable" down

# Install migrate tool
install-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Install development tools
install-tools: install-swag install-migrate
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Full setup for development
setup: install-tools
	go mod tidy
	make swagger
	make docker-up
	sleep 5
	make migrate-up

# Run application with swagger generation
dev: swagger run