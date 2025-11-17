#!/bin/bash

set -e

echo "Starting test database..."
docker-compose -f docker-compose.test.yml up -d

echo "Waiting for database to be ready..."
sleep 10

echo "Running migrations..."
psql "postgres://user:password@localhost:5433/bruschi_rentals_test?sslmode=disable" -c "
-- Create neighborhoods table
CREATE TABLE IF NOT EXISTS neighborhoods (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL
);

-- Create buildings table
CREATE TABLE IF NOT EXISTS buildings (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    neighborhood_id UUID NOT NULL REFERENCES neighborhoods(id) ON DELETE CASCADE,
    address TEXT NOT NULL
);

-- Create index on neighborhood_id for better query performance
CREATE INDEX IF NOT EXISTS idx_buildings_neighborhood_id ON buildings(neighborhood_id);
"

echo "Running e2e tests..."
DATABASE_URL="postgres://user:password@localhost:5433/bruschi_rentals_test?sslmode=disable" go test ./cmd/server/... -v

echo "Cleaning up..."
docker-compose -f docker-compose.test.yml down -v

echo "Tests completed successfully!"
