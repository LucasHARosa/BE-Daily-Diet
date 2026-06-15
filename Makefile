DB_URL=postgres://docker:docker@localhost:5432/daily_diet?sslmode=disable

run:
	go run ./cmd/api

up:
	docker compose up -d

down:
	docker compose down

migrate-up:
	goose -dir db/migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir db/migrations postgres "$(DB_URL)" down

migrate-status:
	goose -dir db/migrations postgres "$(DB_URL)" status

sqlc:
	sqlc generate

test:
	go test ./...

fmt:
	gofmt -w .

tidy:
	go mod tidy

build:
	go build -o bin/api ./cmd/api
