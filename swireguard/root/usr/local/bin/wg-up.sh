#!/usr/bin/with-contenv bash
# shellcheck shell=bash

# This script needs bash. "sh wg-up.sh" ignores the shebang, so fail clearly
if [ -z "${BASH_VERSION:-}" ]
then
    echo "ERROR: $0 must be run with bash, not sh. Use: PostUp = /usr/local/bin/wg-up.sh %i" >&2
    exit 1
fi

# Show connection info once the tunnel is up. Every attempt tries all the
# services in public-ipinfo.sh (a single one, like ipinfo.io, may rate-limit
# shared VPN exit addresses with HTTP 429 and tear down a working tunnel)
#
# Usage: wg-up.sh [interface]   (from wg-quick: PostUp = /usr/local/bin/wg-up.sh %i)
# Without an argument, all the interfaces in "wg show interfaces" are checked.
#
# Then it pings some well known public IPs through the WireGuard interface
# (ping -I binds to it, so the reply proves the tunnel works, not the host
# network). Ping failures are only reported (stderr), they never make the
# script fail: many hosts drop ICMP, so it is not a reliable reason to tear
# down a tunnel.
#  - WG_PING_PUBLIC: space separated public IPs. Default: "1.1.1.1 8.8.8.8 9.9.9.9"
#  - WG_PING_TIMEOUT: seconds to wait for each reply. Default: 2
MAX_ATTEMPTS=60
ATTEMPT=0
WG_UP_INTERFACE="${1:-$(wg show interfaces 2>/dev/null)}"
WG_PING_PUBLIC="${WG_PING_PUBLIC:-1.1.1.1 8.8.8.8 9.9.9.9}"
WG_PING_TIMEOUT="${WG_PING_TIMEOUT:-2}"

# shellcheck source=/dev/null
source /usr/local/bin/public-ipinfo.sh

# Ping the public IPs through every WireGuard interface, report on stderr
ping_checks() {
    local iface host
    local reachable
    local ifaces

    ifaces="${WG_UP_INTERFACE}"
    if [[ -z "${ifaces}" ]]
    then
        echo "WARNING: No WireGuard interface found, skipping ping checks" >&2
        return 0
    fi
    for iface in ${ifaces}
    do
        reachable=0
        for host in ${WG_PING_PUBLIC}
        do
            if ping -I "${iface}" -c 1 -W "${WG_PING_TIMEOUT}" "${host}" >/dev/null 2>&1
            then
                echo "* Ping ${host} via ${iface}: OK" >&2
                reachable=$((reachable + 1))
            else
                echo "WARNING: Ping ${host} via ${iface}: no reply in ${WG_PING_TIMEOUT}s" >&2
            fi
        done
        if [[ -n "${WG_PING_PUBLIC}" ]] && [[ ${reachable} -eq 0 ]]
        then
            echo "WARNING: None of the public IPs [${WG_PING_PUBLIC}] replied via ${iface}" >&2
        fi
    done
}

while [[ ${ATTEMPT} -lt ${MAX_ATTEMPTS} ]]
do
    ATTEMPT=$((ATTEMPT + 1))
    for url in ${IPINFO_SERVICES}
    do
        if OUTPUT=$(ipinfo_query "${url}")
        then
            echo "${OUTPUT}"
            ping_checks
            exit 0
        fi
        echo "WARNING: attempt ${ATTEMPT}/${MAX_ATTEMPTS}: ${url} failed: ${OUTPUT}" >&2
    done
    sleep 1
done

echo "ERROR: Failed to get public IP info from any of [${IPINFO_SERVICES}] after ${MAX_ATTEMPTS} attempts" >&2
ping_checks
exit 1
