#!/bin/bash

# Database migration script
# This script runs the database migrations using Goose

set -e

# Default database configuration
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5454}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-postgres}
DB_SSLMODE=${DB_SSLMODE:-disable}

# Build the database URL
DB_URL="host=$DB_HOST port=$DB_PORT user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME sslmode=$DB_SSLMODE"

echo "Running database migrations..."
echo "Database: $DB_HOST:$DB_PORT/$DB_NAME"
echo "User: $DB_USER"
echo ""

# Run migrations
goose -dir ./migrations postgres "$DB_URL" up

echo ""
echo "Migrations completed successfully!" 