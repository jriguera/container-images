#!/usr/bin/env bash

# Enable strict error handling
set -euo pipefail

# Get the first active interface name
INTERFACE=$(wg show interfaces | awk '{print $1}')
if [[ -z "${INTERFACE}" ]]
then
    echo "* HEALTHCHECK: No active WireGuard interface found" >&2
    exit 1
fi

# Check if the interface is UP using JSON output
# WireGuard interfaces typically show operstate as "UNKNOWN" which is normal
# So we check if the "UP" flag is present in the flags array instead
if ! ip -json link show "${INTERFACE}" 2>/dev/null | jq -e '.[0].flags[] | select(. == "UP")' >/dev/null
then
    echo "* HEALTHCHECK: Interface ${INTERFACE} is not UP" >&2
    exit 1
fi

CURRENT_TIME=$(date +%s)
LATEST_HANDSHAKE=$(wg show "${INTERFACE}" latest-handshakes | awk '{if ($2 > max) max = $2} END {print max}')
if [[ $((CURRENT_TIME - LATEST_HANDSHAKE)) -gt 180 ]]
then
    echo "* HEALTHCHECK: No recent handshakes (> 3 minutes)" >&2
    exit 1
fi

echo "* HEALTHCHECK: WireGuard interface ${INTERFACE} is healthy"
exit 0