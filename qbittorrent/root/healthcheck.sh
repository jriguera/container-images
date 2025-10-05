#!/usr/bin/env bash

# Enable strict error handling
set -euo pipefail

# Set default config directory if not provided
CONFIGDIR="${CONFIGDIR:-/config}"
QBITTORRENT_CONFIG="${CONFIGDIR}/config/qBittorrent.conf"

# Check if qBittorrent.conf exists
if [[ ! -f "${QBITTORRENT_CONFIG}" ]]
then
    echo "* HEALTHCHECK: qBittorrent config file ${QBITTORRENT_CONFIG} not found" >&2
    exit 1
fi

webui_port=$(grep -Po "^WebUI\\\Port=\K(.*)" "${QBITTORRENT_CONFIG}")
webui_address=$(grep -Po "^WebUI\\\Address=\K(.*)" "${QBITTORRENT_CONFIG}")
if [[ -z ${webui_address} ]] || [[ ${webui_address} == "*" ]]
then
    webui_address="127.0.0.1"
fi

if ! nc -z ${webui_address} ${webui_port}
then
    echo "* HEALTHCHECK: qBittorrent WebUI not responding at ${webui_address}:${webui_port}" >&2
    exit 1
fi
