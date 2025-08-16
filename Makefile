.PHONY: install server-dev infra-dev-up infra-dev-down mg-up mg-down mg-reset lint format vuln-check docs help

install: ## Install dependencies and required tools
	go mod download
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.2.2
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install golang.org/x/tools/cmd/gofumpt@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest

server-dev: ## Run server in development mode
	@cd ./internal/cmd/server && go run .

infra-dev-up: ## Start development infrastructure with Docker Compose
	@docker compose --env-file .env.dev -f docker-compose.dev.yml up -d

infra-dev-down: ## Stop and remove development infrastructure
	@docker compose --env-file .env.dev -f docker-compose.dev.yml down -v

mg-up: ## Apply all new database migrations (migrate up)
	@migrate -source file://./migrations -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable up

mg-down: ## Rollback all applied database migrations (migrate down)
	@migrate -source file://./migrations -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable down -all

mg-reset: ## Drop the database schema and reset migrations
	@migrate -source file://./migrations -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable drop -f

lint: ## Run golangci-lint on all Go files
	@golangci-lint run ./...

format: ## Format Go code with gofumpt
	@gofumpt -w -l .

vuln-check: ## Run govulncheck to check for vulnerabilities
	@govulncheck ./...

docs: ## Generate Swagger API docs
	@swag init -g internal/cmd/server/main.go

help: ## Show available make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'