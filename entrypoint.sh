#!/bin/bash
set -e

shutdown() {
    echo "Shutting down gracefully..."
    kill -TERM "$child" 2>/dev/null
    wait "$child"
    exit 0
}

trap shutdown SIGTERM SIGINT

mkdir -p /app/db

python manage.py migrate --noinput

python manage.py collectstatic --noinput --clear 2>/dev/null || true

if [ -n "$RUNSERVER" ]; then
    echo "Starting Django runserver on 0.0.0.0:8000..."
    exec python manage.py runserver 0.0.0.0:8000
fi

export WORKERS=${WORKERS:-3}

echo "Starting gunicorn with $WORKERS workers..."
exec gunicorn chele.wsgi:application \
    --bind 0.0.0.0:8000 \
    --workers "$WORKERS" \
    --timeout 120 \
    --access-logfile - \
    --error-logfile -
