#!/bin/bash
set -e

mkdir -p /app/db

echo "Running migrations..."
python manage.py migrate --noinput

echo "Starting gunicorn..."
exec gunicorn chele.wsgi:application --bind 0.0.0.0:8000 --workers 3