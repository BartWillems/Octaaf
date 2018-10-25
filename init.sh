#!/bin/sh
#
# Startup script used for docker-compose

set -e

go build -o octaaf
./octaaf
