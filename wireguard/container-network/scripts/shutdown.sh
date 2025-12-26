#!/bin/bash

# Clean up routing rules
ip rule del fwmark ${IPTABLES_MANGLE_MARK_PUBLISHED_PORTS} table 200
ip route del default table 200

# Clean up iptables rules
iptables -t nat -D POSTROUTING -s ${INTERNAL_NET_SUBNET} -o wg0 -j MASQUERADE
iptables -D FORWARD -d ${INTERNAL_NET_SUBNET} -i wg0 -m state --state ESTABLISHED,RELATED -j ACCEPT
iptables -D FORWARD -s ${INTERNAL_NET_SUBNET} -o wg0 -j ACCEPT

echo "shutdown done"
