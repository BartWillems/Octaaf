FROM ubuntu:18.04

ENV ENVIRONMENT="development" \
    TELEGRAM_API_KEY="ChangeMe" \
    KALI_ID="0" \
    DATABASE_URI="postgres://octaaf:@127.0.0.1:5432/octaaf_development?sslmode=disable" \
    REDIS_URI="localhost:6379" \
    GOOGLE_API_KEY="ChangeMe" \
    JAEGER_SERVICE_NAME="octaaf" \
    TRUMP_FONT_PATH="/usr/share/fonts/truetype/ubuntu/Ubuntu-LI.ttf"

RUN apt-update \
    && apt install -y fonts-ubuntu \
    && apt clean

RUN mkdir -p /opt/octaaf/config

ADD ./assets /opt/octaaf/assets
ADD ./migrations /opt/octaaf/migrations
ADD ./config/settings.toml.dist /opt/octaaf/config/settings.toml
ADD ./config/database.yml /opt/octaaf/config/database.yml
ADD ./octaaf /opt/octaaf/octaaf

EXPOSE 8080

WORKDIR /opt/octaaf

CMD [ "/opt/octaaf/octaaf" ]