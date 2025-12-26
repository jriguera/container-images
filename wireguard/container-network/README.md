# container-network

A daemon that watches Docker/Podman containers and dynamically manages iptables rules for network routing. It also warms up Linux reverse path routing tables to prevent asymmetric routing issues. Designed to run inside a WireGuard VPN container to provide selective public access to container ports while routing other traffic through an alternative gateway.

## Use Case

This tool solves a specific networking challenge in containerized environments:

```
                                   ┌─────────────────────────────────────────────────┐
                                   │           WireGuard Container                   │
                                   │                                                 │
   Internet ◄───── wg0 ◄───────────┤  container-network daemon                       │
                                   │       │                                         │
                                   │       ▼                                         │
                                   │  ┌─────────────────────────────────────────┐    │
                                   │  │ Watches containers on internal network  │    │
                                   │  │                                         │    │
                                   │  │ • Warm up reverse path routing tables   │    │
                                   │  │ • DNAT ports → route via WireGuard (wg0)│    │
                                   │  │ • Published ports → mark & route via    │    │
                                   │  │   provider network (eth0, table 200)    │    │
                                   │  └─────────────────────────────────────────┘    │
                                   │                                                 │
                                   └────────────────────┬────────────────────────────┘
                                                        │
                                                        │ eth1 (internal network)
                                                        ▼
                                   ┌─────────────────────────────────────────────────┐
                                   │           Internal Docker Network               │
                                   │                                                 │
                                   │  ┌─────────┐  ┌─────────┐  ┌─────────┐          │
                                   │  │Container│  │Container│  │Container│          │
                                   │  │ nginx   │  │  app    │  │   db    │          │
                                   │  │         │  │         │  │         │          │
                                   │  │ Label:  │  │ Publish │  │ No      │          │
                                   │  │ dnat:80 │  │ :8080   │  │ ports   │          │
                                   │  └─────────┘  └─────────┘  └─────────┘          │
                                   │       │             │                           │
                                   └───────┼─────────────┼───────────────────────────┘
                                           │             │
           DNAT via wg0 (VPN) ◄────────────┘             │
           (masquerading)                                │
                                                         │
           Marked packets via eth0 (local network)  ◄────┘
           (routing table 200)
```

### Traffic Routing

1. **DNAT Ports** (labeled containers): Traffic is DNATed to the container and routed through the WireGuard tunnel (`wg0`), appearing to come from the VPN's public IP.

2. **Published Ports** (non-DNAT): Response packets are marked with an iptables mark, which triggers policy routing via an alternative gateway (e.g., provider network `eth0`), bypassing the VPN tunnel.

## Features

- **Container Discovery**: Automatically discovers and monitors running containers on a specified Docker/Podman network
- **Reverse Path Warm-up**: Initiates connections to containers to warm up Linux routing tables, preventing packets from being dropped by reverse path filtering (`rp_filter`)
- **DNAT Rules**: Creates NAT PREROUTING rules to forward external traffic to containers (for VPN-routed public access)
- **Packet Marking**: Marks response packets from published ports for policy-based routing (to bypass VPN)
- **Startup/Shutdown Scripts**: Execute custom scripts for environment setup and cleanup
- **Label-based Filtering**: Only manage containers with a specific label

## Build

```bash
# Clone the repository
git clone https://github.com/yourusername/container-network.git
cd container-network

# Build
make build

# Or build for all platforms
make build-all
```

## Configuration

Configuration can be set via command-line flags or environment variables:

| Flag | Environment Variable | Default | Description |
|------|---------------------|---------|-------------|
| `-runtime-api` | `RUNTIME_API` | auto-detect | Path to Docker/Podman socket |
| `-watch-network` | `WATCH_NETWORK` | `bridge` | Docker network to watch |
| `-watch-container-label` | `WATCH_CONTAINER_LABEL` | `network.enable` | Label that must be `true` on containers |
| `-iptables-mangle-mark-published-ports` | `IPTABLES_MANGLE_MARK_PUBLISHED_PORTS` | (disabled) | Mark value for published port packets |
| `-iptables-dnat-ports-label` | `IPTABLES_DNAT_PORTS_LABEL` | `network.dnat.ports` | Container label specifying DNAT ports |
| `-startup-script` | `STARTUP_SCRIPT` | (none) | Script to run before starting |
| `-shutdown-script` | `SHUTDOWN_SCRIPT` | (none) | Script to run on shutdown |

## Container Labels

Containers must have specific labels to be managed:

| Label | Example Value | Description |
|-------|--------------|-------------|
| `network.enable` | `true` | Enable container watching (required) |
| `network.dnat.ports` | `80,443/tcp,53/udp` | Ports to DNAT (route via VPN) |

## iptables Rules Created

### For DNAT Ports (Public Access via VPN)

When a container has the `network.dnat.ports=80/tcp` label:

```bash
# NAT PREROUTING - redirect incoming traffic to container
iptables -t nat -A PREROUTING -p tcp --dport 80 -j DNAT --to-destination 172.20.0.5:80

# FORWARD - allow traffic to reach the container
iptables -A FORWARD -p tcp -d 172.20.0.5 --dport 80 -j ACCEPT
```

### For Published Ports (Bypass VPN via Mark)

For published ports NOT in the DNAT list, when a container has the published ports

```bash
# MANGLE PREROUTING - mark response packets from published ports
iptables -t mangle -A PREROUTING -p tcp --sport 8080 -j MARK --set-mark 2
```

The mark triggers policy routing via an alternative routing table.

## Example Setup

### WireGuard Container with container-network

**docker-compose.yml:**

```yaml
services:
  wireguard:
    image: your-wireguard-image
    cap_add:
      - NET_ADMIN
      - NET_RAW
    sysctls:
      - net.ipv4.ip_forward=1
      - net.ipv4.conf.all.src_valid_mark=1
    environment:
      - INTERNAL_NET_SUBNET=172.20.0.0/24
      - PROVIDER_NET_GW=192.168.1.1
      - WATCH_NETWORK=internal
      - IPTABLES_MANGLE_MARK_PUBLISHED_PORTS=2
      - STARTUP_SCRIPT=/scripts/startup.sh
      - SHUTDOWN_SCRIPT=/scripts/shutdown.sh
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - ./scripts:/scripts:ro
    networks:
      - internal
      - provider

  nginx:
    image: nginx
    labels:
      - "network.enable=true"
      - "network.dnat.ports=80,443"  # Public via VPN
    ports:
      - "8080:8080"  # Internal access via provider network
    networks:
      - internal

networks:
  internal:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/24
  provider:
    external: true
```

**scripts/startup.sh:**

```bash
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
```

**scripts/shutdown.sh:**

```bash
#!/bin/bash

# Clean up routing rules
ip rule del fwmark ${IPTABLES_MANGLE_MARK_PUBLISHED_PORTS} table 200
ip route del default table 200

# Clean up iptables rules
iptables -t nat -D POSTROUTING -s ${INTERNAL_NET_SUBNET} -o wg0 -j MASQUERADE 2>/dev/null
iptables -D FORWARD -d ${INTERNAL_NET_SUBNET} -i wg0 -m state --state ESTABLISHED,RELATED -j ACCEPT 2>/dev/null
iptables -D FORWARD -s ${INTERNAL_NET_SUBNET} -o wg0 -j ACCEPT 2>/dev/null
```

## How It Works

1. **Startup**: Executes startup script to configure routing tables and base iptables rules

2. **Container Discovery**: Scans for existing containers matching the network and label criteria and keeps monitoring Docker/Podman events for container start/stop

3. **On Container Start**:
   - Warms up reverse path routing by connecting to the container
   - If `network.dnat.ports` label exists: creates DNAT + FORWARD rules
   - For other published ports: creates mangle mark rules

4. **On Container Stop**:
   - Removes all iptables rules created for that container

5. **Shutdown**: Executes shutdown script to clean up routing configuration

## Reverse Path Warm-up

Linux's reverse path filtering (`rp_filter`) can drop packets if the routing path is asymmetric. This daemon "warms up" the routing tables by initiating connections to containers, ensuring the kernel learns the correct routes before traffic flows.

## Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Clean build artifacts
make clean
```

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.
