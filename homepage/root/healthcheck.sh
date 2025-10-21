#!/usr/bin/env bash

# Enable strict error handling
set -euo pipefail

port="${PORT:-3000}"
address="${HOSTNAME:-127.0.0.1}"
[[ "$address" = "0.0.0.0" ]] && address="127.0.0.1"

if ! wget -q --spider http://$address:$port/api/healthcheck
then
    echo "* HEALTHCHECK: Health check failed" >&2
    exit 1
fi
