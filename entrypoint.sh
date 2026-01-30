#!/bin/sh
set -e

DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_NAME}?sslmode=disable"

echo "Running database migrations..."
migrate -database "${DATABASE_URL}" -path ./migrations up

echo "Starting application..."
exec ./main
