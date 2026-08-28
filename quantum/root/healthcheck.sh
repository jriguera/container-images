#!/usr/bin/env bash

# Enable strict error handling
set -euo pipefail

# Set default config directory if not provided
CONFIGDIR="${CONFIGDIR:-/config}"
FB_CONFIG_FILE="${FB_CONFIG_FILE:-${CONFIGDIR}/config.yaml}"

# Check if config file exists
if [[ ! -f "${FB_CONFIG_FILE}" ]]; then
    echo "* HEALTHCHECK: Config file ${FB_CONFIG_FILE} not found" >&2
    exit 1
fi

# Extract listening port and baseURL from the config, with sane defaults.
# On the v2.x line these live under the top-level "http:" key.
port=${FB_PORT:-$(awk '/^http:/{f=1} f&&/^[ \t]+port:/{print $2; exit}' "${FB_CONFIG_FILE}")}
port=${port:-${PORT:-8000}}

base=${FB_BASE_URL:-$(awk '/^http:/{f=1} f&&/^[ \t]+baseURL:/{gsub(/"/,"",$2); print $2; exit}' "${FB_CONFIG_FILE}")}
base=${base:-/}
# normalise: ensure single leading slash, no trailing slash (unless root)
base="/${base#/}"
base="${base%/}"

if ! curl -fsS "http://127.0.0.1:${port}${base}/health" >/dev/null; then
    echo "* HEALTHCHECK: Health check failed on http://127.0.0.1:${port}${base}/health" >&2
    exit 1
fi
