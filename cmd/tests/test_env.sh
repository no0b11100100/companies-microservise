#!/usr/bin/env bash

# Database connection settings
export DB_HOST="${DB_HOST:-localhost}"
export DB_PORT="${DB_PORT:-3306}"
export DB_NAME="${DB_NAME:-companiesdb}"
export DB_USER="${DB_USER:-root}"
export DB_PASSWORD="${DB_PASSWORD:-password}"

# Kafka settings
export KAFKA_BROKER="${KAFKA_BROKER:-localhost:9092}"
export KAFKA_TOPIC="${KAFKA_TOPIC:-company_events}"

# HTTP server settings
export HTTP_HOST="${HTTP_HOST:-0.0.0.0}"
export HTTP_PORT="${HTTP_PORT:-8080}"

# Optional: Logging and auth
export JWT_SECRET="${JWT_SECRET:-mysecretkey}"

echo "Environment variables set:"
echo "  DB_HOST=$DB_HOST"
echo "  DB_PORT=$DB_PORT"
echo "  DB_NAME=$DB_NAME"
echo "  DB_USER=$DB_USER"
echo "  DB_PASSWORD=test"
echo "  KAFKA_BROKER=$KAFKA_BROKER"
echo "  HTTP_HOST=$HTTP_HOST"
echo "  HTTP_PORT=$HTTP_PORT"
echo "  JWT_SECRET=w4X8bY9s4fJ7nHxPlA+2Fv7Yq8RmQO3SZj8Phq+k6Xo="
