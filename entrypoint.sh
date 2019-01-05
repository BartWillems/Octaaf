#!/bin/bash

set -e

if [ "$TZ" != $(cat /etc/timezone) ];
then
    echo "$TZ" > /etc/timezone
    ln -snf "/usr/share/zoneinfo/$TZ" /etc/localtime
    dpkg-reconfigure -f noninteractive tzdata
fi

exec "$@"