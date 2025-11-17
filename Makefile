.PHONY: build run dev clean test lint coverage release docker-dev docker-prod migrate-up migrate-down migrate-create

DB_URL ?= $(shell echo $$DATABASE_URL || echo postgres://user:password@localhost:5432/bruschi_rentals?sslmode=disable)

build:
	go build -o server ./cmd/server

run: build
	./server

dev:
	air

test:
	go test ./cmd/server/... -v

lint:
	golangci-lint run

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

release:
	@echo "Usage: make release VERSION=1.0.0"
	git tag v$(VERSION)
	git push origin v$(VERSION)

migrate-up:
	migrate -path ./migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DB_URL)" down

migrate-force:
	migrate -path ./migrations -database "$(DB_URL)" force 1

migrate-create:
	@echo "Usage: make migrate-create NAME=migration_name"

docker-dev:
	docker-compose -f docker-compose.dev.yml up --build

clean:
	rm -f server coverage.out
	docker-compose down -v

clean-all: clean
	rm -rf tmp/ postgres_data/