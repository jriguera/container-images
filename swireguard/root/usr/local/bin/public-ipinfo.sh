#!/usr/bin/bash
# shellcheck shell=bash

# Print public IP information as JSON (stdout). Several services are tried in
# order, because they rate-limit shared VPN exit addresses (HTTP 429). Failures
# are reported on stderr, so stdout only carries the JSON document.
# Override the list with IPINFO_SERVICES (space separated URLs).

IPINFO_SERVICES="${IPINFO_SERVICES:-https://ipwho.is/ https://ipinfo.io https://ipapi.co/json}"
IPINFO_TIMEOUT="${IPINFO_TIMEOUT:-5}"

# Query one service. Prints the JSON on success, the reason on failure
ipinfo_query() {
    local url="$1"
    local output

    if ! output=$(wget -q -T "${IPINFO_TIMEOUT}" -t 1 -O - "${url}" 2>&1)
    then
        echo "${output:-connection failed, timeout or HTTP error}"
        return 1
    fi
    if [[ -z "${output}" ]]
    then
        echo "empty response"
        return 1
    fi
    if [[ "${output}" != *"{"* ]]
    then
        echo "not a JSON response: ${output:0:200}"
        return 1
    fi
    # Some services answer HTTP 200 with an error document (ipwho.is, ipapi.co)
    if [[ "${output}" =~ \"success\"[[:space:]]*:[[:space:]]*false ]] || \
       [[ "${output}" =~ \"error\"[[:space:]]*:[[:space:]]*true ]]
    then
        echo "service error: $(tr -d '\n' <<< "${output}" | cut -c1-200)"
        return 1
    fi
    echo "${output}"
}

ipinfo_json() {
    local url
    local output

    for url in ${IPINFO_SERVICES}
    do
        if output=$(ipinfo_query "${url}")
        then
            printf '%s\n' "${output}"
            return 0
        fi
        echo "WARNING: ${url} failed: ${output}" >&2
    done
    printf '{\n'
    printf '\t"error": "Failed to retrieve IP information",\n'
    printf '\t"message": "All services failed: %s",\n' "${IPINFO_SERVICES}"
    printf '\t"timeout": %d\n' "${IPINFO_TIMEOUT}"
    printf '}\n'
    return 1
}

# Only run when executed, so other scripts can source the functions
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]
then
    ipinfo_json
fi
