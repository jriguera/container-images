#!/usr/bin/env bash

# Enable strict error handling
set -euo pipefail

if ! curl --fail http://127.0.0.1:${PORT}
then
    echo "* HEALTHCHECK: Health check failed" >&2
    exit 1
fi
