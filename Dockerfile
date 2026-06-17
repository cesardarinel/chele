FROM python:3.13-slim AS builder

ENV PYTHONDONTWRITEBYTECODE=1

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc libpq-dev && \
    rm -rf /var/lib/apt/lists/*

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt


FROM python:3.13-slim

ENV \
    PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    APP_HOME=/app \
    APP_USER=appuser \
    APP_USER_UID=1001 \
    HOME=/tmp

RUN addgroup --system --gid $APP_USER_UID $APP_USER \
    && adduser --system --uid $APP_USER_UID --gid $APP_USER_UID $APP_USER \
    && apt-get update && apt-get install -y --no-install-recommends \
        libpq5 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR $APP_HOME

COPY --from=builder /usr/local/lib/python3.13/site-packages /usr/local/lib/python3.13/site-packages
COPY --from=builder /usr/local/bin /usr/local/bin

COPY . .

RUN chmod +x $APP_HOME/entrypoint.sh \
    && mkdir -p $APP_HOME/staticfiles $APP_HOME/db \
    && python manage.py collectstatic --noinput \
    && chown -R $APP_USER:$APP_USER $APP_HOME

USER $APP_USER

ENTRYPOINT ["/app/entrypoint.sh"]
