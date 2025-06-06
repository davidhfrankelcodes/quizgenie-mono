#!/bin/sh
# wait-for-postgres.sh 

set -e

# Usage:   ./wait-for-postgres.sh <your_binary> [args...]
# Example: ./wait-for-postgres.sh quizgenie_api

# Loop until pg_isready returns “accepting connections”
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" >/dev/null 2>&1; do
  echo "Waiting for Postgres at $DB_HOST:$DB_PORT..."
  sleep 2
done

echo "Postgres is up — starting the API."
exec "$@"
