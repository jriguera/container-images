#!/usr/bin/env bash

# Enable strict error handling
set -euo pipefail

# Set default config directory if not provided
CONFIGDIR="${CONFIGDIR:-/config}"

# Check if settings.json exists
if [[ ! -f "${CONFIGDIR}/settings.json" ]]; then
    echo "* HEALTHCHECK: Settings file ${CONFIGDIR}/settings.json not found" >&2
    exit 1
fi

# Extract port and address from settings.json, with defaults
port=${FB_PORT:-$(jq -r .port "${CONFIGDIR}/settings.json")}
address=${FB_ADDRESS:-$(jq -r .address "${CONFIGDIR}/settings.json")}
address=${address:-127.0.0.1}

if ! wget -q --spider http://$address:$port/health
then
    echo "* HEALTHCHECK: Health check failed" >&2
    exit 1
fi