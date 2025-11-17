#!/bin/bash

set -e

echo "Starting test database..."
docker-compose -f docker-compose.test.yml up -d

echo "Waiting for database to be ready..."
sleep 10

echo "Running migrations..."
make migrate-up DB_URL="postgres://user:password@localhost:5433/bruschi_rentals_test?sslmode=disable"

echo "Running e2e tests..."
DATABASE_URL="postgres://user:password@localhost:5433/bruschi_rentals_test?sslmode=disable" go test ./cmd/server/... -v

echo "Cleaning up..."
docker-compose -f docker-compose.test.yml down -v

echo "Tests completed successfully!"
