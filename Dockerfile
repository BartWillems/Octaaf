FROM ubuntu:18.04

ENV ENVIRONMENT="development" \
    TELEGRAM_API_KEY="ChangeMe" \
    KALI_ID="0" \
    DATABASE_URI="postgres://octaaf:@127.0.0.1:5432/octaaf_development?sslmode=disable" \
    REDIS_URI="localhost:6379" \
    REDIS_DB="0" \
    GOOGLE_API_KEY="ChangeMe" \
    JAEGER_SERVICE_NAME="octaaf" \
    JAEGER_AGENT_HOST="localhost" \
    JAEGER_AGENT_PORT="6831" \
    TRUMP_FONT_PATH="/usr/share/fonts/truetype/ubuntu/Ubuntu-LI.ttf"

RUN apt update \
    && apt install -y --no-install-recommends fonts-ubuntu ca-certificates \
    && apt clean \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /opt/octaaf/config

ADD ./assets /opt/octaaf/assets
ADD ./migrations /opt/octaaf/migrations
ADD ./config/settings.toml.dist /opt/octaaf/config/settings.toml
ADD ./config/database.yml /opt/octaaf/config/database.yml
ADD ./octaaf /opt/octaaf/octaaf

# Production port - development port
EXPOSE 8080 8888

WORKDIR /opt/octaaf

CMD [ "/opt/octaaf/octaaf" ]