.PHONY: install server-dev infra-dev-up infra-dev-down mg-up mg-down lint

install:
	go mod download
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.2.2
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

server-dev:
	@cd ./internal/cmd/server && go run .

infra-dev-up:
	@docker compose --env-file .env.dev -f docker-compose.dev.yml up -d

infra-dev-down:
	@docker compose --env-file .env.dev -f docker-compose.dev.yml down -v

mg-up:
	@migrate -source file://./migrations -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable up

mg-down:
	@migrate -source file://./migrations -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable down -all

mg-reset:
	@migrate -source file://./migrations -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable drop -f

lint:
	@golangci-lint run ./...