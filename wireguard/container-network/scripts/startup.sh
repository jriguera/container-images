#!/bin/bash

# Disable reverse path filtering (required for asymmetric routing)
sysctl -w net.ipv4.conf.wg0.rp_filter=0
sysctl -w net.ipv4.conf.eth1.rp_filter=0

# Allow forwarding between internal network and WireGuard
iptables -I FORWARD -s ${INTERNAL_NET_SUBNET} -o wg0 -j ACCEPT
iptables -I FORWARD -d ${INTERNAL_NET_SUBNET} -i wg0 -m state --state ESTABLISHED,RELATED -j ACCEPT
# Masquerade internal traffic going out via WireGuard
iptables -t nat -I POSTROUTING 1 -s ${INTERNAL_NET_SUBNET} -o wg0 -j MASQUERADE

# Policy routing: marked packets (fwmark 2) use table 200 (provider gateway)
ip route add default via ${PROVIDER_NET_GW} table 200
ip rule add fwmark ${IPTABLES_MANGLE_MARK_PUBLISHED_PORTS} table 200

echo "startup done"
