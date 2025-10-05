#!/usr/bin/env bash

# Enable strict error handling
set -euo pipefail

# Set default config directory if not provided
CONFIGDIR="${CONFIGDIR:-/config}"

HEALTHCHECK_CLIENT_ID=healthcheck
HEALTHCHECK_TOPIC='$SYS/broker/uptime'

# Check if mosquitto.conf exists
if [[ ! -f "${CONFIGDIR}/mosquitto.conf" ]]; then
    echo "* HEALTHCHECK: Settings file ${CONFIGDIR}/mosquitto.conf not found" >&2
    exit 1
fi

# find a listener using sockets
if SOCKET=$(sed -n '/^listener\s\+0 .*/{s/^listener\s\+0\s\+\(.*\)$/\1/p;h};${x;/./{x;q0};x;q1}' "${CONFIGDIR}/mosquitto.conf")
then
    mosquitto_sub --unix ${SOCKET} -t ${HEALTHCHECK_TOPIC} -E -C 1 -i ${HEALTHCHECK_CLIENT_ID} -W 3
else
    mosquitto_sub -h 127.0.0.1 -t ${HEALTHCHECK_TOPIC} -E -C 1 -i ${HEALTHCHECK_CLIENT_ID} -W 3
fi