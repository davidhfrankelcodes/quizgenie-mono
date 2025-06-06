#!/bin/sh
set -e

# wait-for-postgres.sh
# Usage: wait-for-postgres.sh <your_binary> [args...]

until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" >/dev/null 2>&1; do
  echo "Waiting for Postgres at $DB_HOST:$DB_PORT..."
  sleep 2
done

echo "Postgres is up—starting the process."
exec "$@"
