#!/usr/bin/env bash

# Enable strict error handling
set -euo pipefail

# Set default config directory if not provided
CONFIGDIR="${CONFIGDIR:-/config}"

# Check if config.xml exists
if [[ ! -f "${CONFIGDIR}/config.xml" ]]; then
    echo "* HEALTHCHECK: Settings file ${CONFIGDIR}/config.xml not found" >&2
    exit 1
fi

PORT=$(xmlstarlet sel -T -t -v /Config/Port ${CONFIGDIR}/config.xml)

if [[ $(curl -sL "http://127.0.0.1:${PORT:-9696}/ping" | jq -r '.status' 2>/dev/null) = "OK" ]]
then
    exit 0
else
    echo "* HEALTHCHECK: Health check failed" >&2
    exit 1
fi