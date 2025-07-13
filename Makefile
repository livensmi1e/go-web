.PHONY: server-dev, db-dev-up, db-dev-down, mg-up, mg-down

server-dev:
	@go run .

db-dev-up:
	@docker compose --env-file .env.development -f docker-compose.dev.yml up -d

db-dev-down:
	@docker compose --env-file .env.development -f docker-compose.dev.yml down -v

mg-up:
	@migrate -source file://./migrations -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable up

mg-down:
	@migrate -source file://./migrations -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable down -all

mg-reset:
	@migrate -source file://./migrations -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable drop -f