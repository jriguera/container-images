#!/usr/bin/with-contenv bash
# shellcheck shell=bash

# Function to get IP information as JSON
ipinfo_json() {
    local url="$1"
    local timeout=5

    # Try to fetch IP info with timeout
    if ! output=$(wget -q --timeout="$timeout" --tries=1 "$url" -O - 2>&1); then
        printf '{\n'
        printf '\t"error": "Failed to retrieve IP information",\n'
        printf '\t"message": "Unable to connect to ipinfo.io or DNS resolution failed",\n'
        printf '\t"timeout": %d\n' "$timeout"
        printf '}\n'
        return 1
    fi
    printf '%s\n' "$output"
    return 0
}

ipinfo_json "https://ipinfo.io"
