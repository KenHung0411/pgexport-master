#!/bin/sh

set -e

cat /MANIFEST

echo "INFO: Waiting for Postgres to start..."
while ! nc -z ${DATABASE_HOST} ${DATABASE_PORT}; do sleep 0.1; done
echo "INFO: Postgres is up"

exec "$@"
