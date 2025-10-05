#!/usr/bin/env bash

# Enable strict error handling
set -euo pipefail

# Set default config directory if not provided
CONFIGDIR="${CONFIGDIR:-/config}"

# Check if remote.conf exists
if [[ ! -f "${CONFIGDIR}/remote.conf" ]]
then
    echo "* HEALTHCHECK: Remote config file ${CONFIGDIR}/remote.conf not found" >&2
    exit 1
fi

cmd="amulecmd --command status"
# Execute amulecmd and capture both stdout and exit code
if output=$(HOME=${CONFIGDIR} ${cmd} 2>&1)
then
    # Check if the output contains "Succeeded!"
    if echo "$output" | grep -q "Succeeded!"
    then
        echo "* HEALTHCHECK: Amule is healthy"
        exit 0
    else
        echo "* HEALTHCHECK: Amule status check did not return success indicator" >&2
        exit 1
    fi
else
    # amulecmd failed
    exit_code=$?
    echo "* HEALTHCHECK: amulecmd failed with exit code $exit_code" >&2
    echo "* HEALTHCHECK: Output: $output" >&2
    exit 1
fi
