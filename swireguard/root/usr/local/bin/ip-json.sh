#!/usr/bin/with-contenv bash
# shellcheck shell=bash

# Function to get network interface information as JSON
iface_to_json() {
    local interface="${1}"
    
    if [[ -z "$interface" ]]; then
        printf '{\n'
        printf '\t"error": "No interface specified",\n'
        printf '\t"message": "Please provide a network interface name"\n'
        printf '}\n'
        return 1
    fi
    
    # Check if interface exists
    if ! ip link show "$interface" &>/dev/null; then
        printf '{\n'
        printf '\t"error": "Interface not found",\n'
        printf '\t"message": "Interface %s does not exist"\n' "$interface"
        printf '}\n'
        return 1
    fi
    
    # Get address info and statistics
    local addr_json=$(ip -json addr show "$interface" 2>/dev/null)
    local stats_json=$(ip -json -s link show "$interface" 2>/dev/null)

    # Extract stats from stats_json and merge it into addr_json
    echo "$addr_json" | jq --argjson stats "$stats_json" '.[0] + {stats: $stats[0].stats64}'
}

# Main execution
iface_to_json "$1"