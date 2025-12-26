#!/usr/bin/with-contenv bash
# shellcheck shell=bash

# Show connection info
MAX_ATTEMPTS=60
ATTEMPT=0

while [[ ${ATTEMPT} -lt ${MAX_ATTEMPTS} ]]; do
    ATTEMPT=$((ATTEMPT + 1))
    if OUTPUT=$(wget -q https://ipinfo.io -O - 2>/dev/null) && [[ -n "${OUTPUT}" ]]; then
        echo "${OUTPUT}"
        exit 0
    fi
    sleep 1
done

echo "ERROR: Failed to reach ipinfo.io after ${MAX_ATTEMPTS} attempts"
exit 1