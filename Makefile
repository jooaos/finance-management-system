include .env

DB_CONNECTION="mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?charset=utf8"

air-up:
	@air -c cmd/${SERVICE}/.air.toml

down:
	@docker compose down

mocks:
	@go generate ./...

migration-create:
	@docker compose --profile tools run --rm migrate ${DB_CONNECTION} create -ext sql -seq -dir /migrations/ ${NAME}

migration-down:
	@docker compose --profile tools run --rm migrate ${DB_CONNECTION} down

migration-up:
	@docker compose --profile tools run --rm migrate ${DB_CONNECTION} up

up:
	@docker compose up

up-dettached:
	@docker compose up -d

up-service:
	@docker compose up ${SERVICE} -d

set-up:
	@go mod download

test-unit:
	@go test ./...