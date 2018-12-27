FROM ubuntu:18.04

RUN mkdir -p /opt/octaaf/config

ADD ./assets /opt/octaaf/assets
ADD ./migrations /opt/octaaf/migrations
ADD ./config/settings.toml.dist /opt/octaaf/config/settings.toml
ADD ./config/database.yml /opt/octaaf/config/database.yml
ADD ./octaaf /opt/octaaf/octaaf

EXPOSE 8080

WORKDIR /opt/octaaf

CMD [ "/opt/octaaf/octaaf" ]