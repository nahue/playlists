#!/bin/bash

# Database connection test script

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

echo "Testing database connection..."
echo "Database: $DB_HOST:$DB_PORT/$DB_NAME"
echo "User: $DB_USER"
echo ""

# Test connection using psql
if command -v psql &> /dev/null; then
    echo "Testing with psql..."
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT version();" || {
        echo "❌ Database connection failed with psql"
        exit 1
    }
    echo "✅ Database connection successful with psql"
else
    echo "⚠️  psql not found, skipping psql test"
fi

echo ""

# Test connection using Goose
echo "Testing with Goose..."
goose -dir ./migrations postgres "$DB_URL" status || {
    echo "❌ Database connection failed with Goose"
    exit 1
}
echo "✅ Database connection successful with Goose"

echo ""
echo "🎉 All database connection tests passed!" 