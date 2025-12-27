# docker-wireguard

WireGuard Docker image based on Alpine Linux with s6-overlay, multi-arch support.

WireGuard is a modern, fast, and secure VPN tunnel that uses state-of-the-art cryptography.

Most of this work is based on https://github.com/linuxserver/docker-wireguard

## Features

- **Multi-arch support**: AMD64, ARM64, ARMv7
- **s6-overlay**: Proper process supervision and initialization
- **Server mode**: Automatic server and peer configuration generation
- **QR code generation**: Easy mobile client setup
- **CoreDNS integration**: Built-in DNS server for VPN clients
- **SOCKS5 proxy**: Optional proxy support
- **Webhook API**: HTTP endpoints for management and monitoring
- **Health checks**: Built-in container health monitoring
- **Container Network**: Dynamic iptables management for routing container traffic via VPN or provider network

## Quick Start

### Basic Server Setup

```bash
docker run -d \
  --name=wireguard \
  --cap-add=NET_ADMIN \
  --cap-add=SYS_MODULE \
  -e PUID=1000 \
  -e PGID=1000 \
  -e TZ=Etc/UTC \
  -e WG_PEERS=3 \
  -p 51820:51820/udp \
  -v /path/to/config:/config \
  -v /lib/modules:/lib/modules:ro \
  --sysctl="net.ipv4.conf.all.src_valid_mark=1" \
  --restart unless-stopped \
  ghcr.io/jriguera/wireguard:latest
```

### Using Docker Compose

See [docker-compose.yml](docker-compose.yml) for a complete example.

```bash
docker-compose up -d
```

## Environment Variables

### Core Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PUID` | `1000` | User ID for running the container |
| `PGID` | `1000` | Group ID for running the container |
| `TZ` | `Etc/UTC` | Timezone (e.g., `America/New_York`, `Europe/London`) |

### WireGuard Server Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `WG_PEERS` | - | **Required for server mode.** Number of peers (e.g., `3`) or comma-separated names (e.g., `laptop,phone,tablet`) |
| `WG_SERVER_PORT` | `51820` | WireGuard server port (UDP) |
| `WG_SERVER_ADDR` | _auto-detect_ | External server address/hostname. Auto-detects public IP if not set |
| `WG_INTERNAL_SUBNET` | `10.13.13.0` | VPN internal subnet (first 3 octets, e.g., `10.13.13`) |
| `WG_ALLOWEDIPS` | `0.0.0.0/0, ::/0` | Traffic routing for peers. Use `0.0.0.0/0` for all traffic through VPN |
| `WG_PEERS_DNS` | _auto (gateway)_ | DNS servers for peers (e.g., `1.1.1.1, 1.0.0.1`). Defaults to VPN gateway (.1 address) |
| `WG_SHOW_QR` | `true` | Display QR codes in container logs. Set to `false` to only save to files |
| `WG_PEERS_PERSISTENTKEEPALIVE` | - | Peers requiring keepalive: `all` or comma-separated names (e.g., `laptop,phone`) |

### Per-Peer Configuration

You can customize AllowedIPs for specific peers using environment variables:

```bash
WG_SERVER_ALLOWEDIPS_PEER_laptop="192.168.1.0/24, 10.0.0.0/8"
WG_SERVER_ALLOWEDIPS_PEER_phone="10.13.13.0/24"
```

Format: `WG_SERVER_ALLOWEDIPS_PEER_<peername>=<cidrs>`

### Webhook API Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `GENERATE_WEBHOOK_CONFIG` | `true` | Auto-generate default webhook config if missing |
| `WEBHOOK_PORT` | `9000` | HTTP port for webhook API |
| `WEBHOOK_IP` | `0.0.0.0` | IP address to bind webhook server |
| `WEBHOOK_VERBOSE` | `false` | Enable verbose logging for webhook requests |
| `WEBHOOK_HOTRELOAD` | `true` | Enable hot-reloading of webhook configuration |

### Other Services

| Variable | Default | Description |
|----------|---------|-------------|
| `COREDNS_PORT` | `53` | CoreDNS server port |
| `PROXY_PORT` | `1080` | SOCKS5 proxy port (if enabled) |

### Default Gateway Priority Configuration

When running WireGuard in a container with an existing default route and different networks, you may want changing the default gw before WG starts. The container can automatically adjust the default gateway's metric (priority) at startup.

| Variable | Default | Description |
|----------|---------|-------------|
| `PRESET_DEFAULT_GW_IP` | _(none)_ | **Required to enable.** IP address of the existing default gateway |
| `PRESET_DEFAULT_GW_DEV` | _(none)_ | Network device for the gateway (e.g., `eth0`). Optional |
| `PRESET_DEFAULT_GW_PRIO` | `0` | Metric (priority) for the new default route. Higher values = lower priority |

When `PRESET_DEFAULT_GW_IP` is set, the container will:

1. Add a new default route with the specified metric (low priority)
2. Delete the original default route (metric 0)

This allows WireGuard to become the preferred route while keeping the provider network as a fallback.

**Example:**

```bash
docker run -d \
  --name=wireguard \
  --cap-add=NET_ADMIN \
  -e PRESET_DEFAULT_GW_IP=192.168.1.1 \
  -e PRESET_DEFAULT_GW_DEV=eth0 \
  -e PRESET_DEFAULT_GW_PRIO=500 \
  ...
```

### Container Network Configuration

The `container-network` daemon watches Docker/Podman containers and manages iptables rules for routing. It enables selective public access to container ports via the VPN while routing other traffic through the provider network. See [container-network/README.md](container-network/README.md) for detailed documentation.

| Variable | Default | Description |
|----------|---------|-------------|
| `RUNTIME_API` | _auto-detect_ | Path to Docker/Podman socket |
| `WATCH_NETWORK` | `bridge` | Docker network to watch for containers |
| `WATCH_CONTAINER_LABEL` | `network.enable` | Label that must be `true` on containers to be managed |
| `IPTABLES_MANGLE_MARK_PUBLISHED_PORTS` | `2` | Mark value for published port packets |
| `IPTABLES_DNAT_PORTS_LABEL` | `network.dnat.ports` | Container label specifying ports to DNAT |
| `STARTUP_SCRIPT` | `/usr/local/bin/container-network-startup.sh` | Script to run before starting |
| `SHUTDOWN_SCRIPT` | `/usr/local/bin/container-network-shutdown.sh` | Script to run on shutdown |

For the default startup and shutdown scripts these environment variables are needed:

| Variable | Default | Description |
|----------|---------|-------------|
| `INTERNAL_NET_SUBNET` | _(none)_ | Internal subnet |
| `INTERNAL_NET_GW` | _(none)_ | Gateway in the internal subnet |

## Volume Mounts

| Path | Description |
|------|-------------|
| `/config` | Configuration directory containing WireGuard configs, keys, and peer data |
| `/lib/modules` | **Required (read-only).** Kernel modules for WireGuard |

### Configuration Directory Structure

```
/config/
├── wg_confs/           # WireGuard interface configs
│   └── wg0.conf
├── server/             # Server keys
│   ├── privatekey-server
│   └── publickey-server
├── templates/          # Configuration templates (customizable)
│   ├── server.conf
│   └── peer.conf
├── peer1/              # Peer configurations and keys
│   ├── peer1.conf
│   ├── peer1.png       # QR code
│   ├── privatekey-peer1
│   ├── publickey-peer1
│   └── presharedkey-peer1
├── peer2/
│   └── ...
└── webhook/            # Webhook configuration
    └── hooks.yaml
```

## Ports

| Port | Protocol | Description |
|------|----------|-------------|
| `51820` | UDP | WireGuard VPN (customizable via `WG_SERVER_PORT`) |
| `9000` | TCP | Webhook API (customizable via `WEBHOOK_PORT`) |
| `53` | UDP | CoreDNS server (optional) |
| `1080` | TCP | SOCKS5 proxy (optional) |

## Webhook API Endpoints

The container includes a webhook API for management and monitoring. See [WEBHOOK_CONFIG.md](WEBHOOK_CONFIG.md) for complete documentation.

### Available Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /wireguard-status` | Get WireGuard interface and peer status (JSON) |
| `GET /network-interface-info?interface=<name>` | Get network interface details (JSON) |
| `GET /public-ip-info` | Get public IP information (JSON) |
| `GET /show-peer-qr?peer=<id>` | Generate peer QR code (text) |

### Example API Calls

```bash
# Get WireGuard status
curl "http://localhost:9000/wireguard-status"

# Get interface information
curl "http://localhost:9000/network-interface-info?interface=wg0"

# Get public IP
curl "http://localhost:9000/public-ip-info"

# Show peer QR code
curl "http://localhost:9000/show-peer-qr?peer=1"
```

## Usage Examples

### Server Mode with Named Peers

```bash
docker run -d \
  --name=wireguard \
  --cap-add=NET_ADMIN \
  --cap-add=SYS_MODULE \
  -e WG_PEERS=laptop,phone,tablet \
  -e WG_SERVER_ADDR=vpn.example.com \
  -e WG_PEERS_DNS=1.1.1.1,1.0.0.1 \
  -p 51820:51820/udp \
  -p 9000:9000 \
  -v ./config:/config \
  -v /lib/modules:/lib/modules:ro \
  --sysctl="net.ipv4.conf.all.src_valid_mark=1" \
  wireguard:latest
```

### Split Tunnel Configuration

Route only specific networks through VPN:

```bash
docker run -d \
  --name=wireguard \
  --cap-add=NET_ADMIN \
  -e WG_PEERS=3 \
  -e WG_ALLOWEDIPS="10.0.0.0/8, 192.168.0.0/16" \
  -p 51820:51820/udp \
  -v ./config:/config \
  wireguard:latest
```

### Custom Per-Peer AllowedIPs

```bash
docker run -d \
  --name=wireguard \
  --cap-add=NET_ADMIN \
  -e WG_PEERS=laptop,phone \
  -e WG_SERVER_ALLOWEDIPS_PEER_laptop="192.168.1.0/24" \
  -e WG_SERVER_ALLOWEDIPS_PEER_phone="10.13.13.0/24" \
  -p 51820:51820/udp \
  -v ./config:/config \
  wireguard:latest
```

### With Webhook API and Verbose Logging

```bash
docker run -d \
  --name=wireguard \
  --cap-add=NET_ADMIN \
  -e WG_PEERS=3 \
  -e WEBHOOK_PORT=9000 \
  -e WEBHOOK_VERBOSE=true \
  -p 51820:51820/udp \
  -p 9000:9000 \
  -v ./config:/config \
  wireguard:latest
```

## Configuration Regeneration

The container automatically regenerates server and peer configurations when these environment variables change:

- `WG_SERVER_ADDR`
- `WG_SERVER_PORT`
- `WG_INTERNAL_SUBNET`
- `WG_PEERS_DNS`
- `WG_ALLOWEDIPS`
- `WG_PEERS_PERSISTENTKEEPALIVE`

**Note**: Existing keys are preserved during regeneration. To regenerate keys, delete the corresponding peer folders.

## Customizing Configuration Templates

You can customize the WireGuard configuration templates:

1. Start the container to generate default templates
2. Edit templates in `/config/templates/`:
   - `server.conf` - Server configuration template
   - `peer.conf` - Peer configuration template
3. Restart the container to apply changes

Templates support variable substitution using bash syntax (e.g., `${WG_SERVER_PORT}`).

## Health Checks

The container includes built-in health checks that monitor:

- WireGuard interface status
- Network connectivity
- Service availability

Check container health:

```bash
docker inspect --format='{{.State.Health.Status}}' wireguard
```

## Security Considerations

### Webhook API

The webhook API exposes management endpoints. **Secure it properly:**

1. **Bind to localhost only** (for local access):
   ```bash
   -e WEBHOOK_IP=127.0.0.1
   ```

2. **Use a reverse proxy** with authentication (nginx, Traefik, etc.)

3. **Add token authentication** to `hooks.yaml`:
   ```yaml
   trigger-rule:
     match:
       type: "value"
       value: "your-secret-token"
       parameter:
         source: "header"
         name: "X-Webhook-Token"
   ```

4. **Enable HTTPS** via reverse proxy

See [WEBHOOK_CONFIG.md](WEBHOOK_CONFIG.md) for detailed security configuration.

## Building the Image

```bash
docker build \
  --no-cache \
  --pull \
  --build-arg BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
  --build-arg VERSION=1.0.0 \
  --build-arg COMMIT_SHA=$(git rev-parse HEAD) \
  -t wireguard:latest \
  .
```

## Troubleshooting

### WireGuard interface not starting

1. Check kernel module availability:
   ```bash
   docker exec wireguard modprobe wireguard
   ```

2. Ensure proper capabilities:
   ```bash
   --cap-add=NET_ADMIN --cap-add=SYS_MODULE
   ```

3. Check sysctl setting:
   ```bash
   --sysctl="net.ipv4.conf.all.src_valid_mark=1"
   ```

### Peers can't connect

1. Verify port forwarding: Ensure UDP port `51820` (or your custom port) is forwarded
2. Check server address: Verify `WG_SERVER_ADDR` is correct and accessible
3. Review logs: `docker logs wireguard`

### Webhook API not responding

1. Check if service is running:
   ```bash
   docker exec wireguard nc -z localhost 9000
   ```

2. Enable verbose logging:
   ```bash
   -e WEBHOOK_VERBOSE=true
   ```

3. Verify configuration:
   ```bash
   docker exec wireguard cat /config/webhook/hooks.yaml
   ```

## Additional Documentation

- [WEBHOOK_CONFIG.md](WEBHOOK_CONFIG.md) - Webhook API detailed documentation
- [WireGuard Official Documentation](https://www.wireguard.com/)
- [s6-overlay Documentation](https://github.com/just-containers/s6-overlay)
