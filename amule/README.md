# docker-amule

Amule Docker image based on Alpine, multi-arch.
Amule is a free P2P file sharing application that works with the eDonkey2000 network and Kademlia network.

### Develop and test builds

Just type:

```
docker build \
    --no-cache \
    --pull \
    --build-arg BUILD_DATE=$(date '+%Y-%m-%dT%H:%M:%S%:z') \
    --build-arg VERSION=0.1 \
    --build-arg COMMIT_SHA=main \
    .  -t amule
```

# Usage

Given the docker image with name `amule`:

```bash
docker run --rm -ti --name amule \
  -p 4662:4662 \
  -p 4672:4672/udp \
  -p 4711:4711 \
  -p 4712:4712 \
  -v $(pwd)/data:/data \
  -v $(pwd)/config:/config \
  amule
```

Or using docker-compose:

```yaml
version: "3.9"
services:
  amule:
    image: amule:latest
    restart: unless-stopped
    ports:
      - "4662:4662"      # eMule TCP
      - "4672:4672/udp"  # eMule UDP  
      - "4711:4711"      # Web interface
      - "4712:4712"      # External Connect
    volumes:
      - ./data:/data
      - ./config:/config
    environment:
      - AMULE_PASSWORD=mysecretpassword
      - AMULE_WEBSERVER_TEMPLATE=adaptable
      - RECURSIVE_SHARE=true
```

You can access the web interface at `http://localhost:4711` using the configured password.

## Environment Variables

### Core Configuration

```bash
# Base directories
DATADIR=/data                     # Data directory (default: /data)
CONFIGDIR=/config                # Config directory (default: /config)
PORT=4711                        # Default port (default: 4711)
AMULE_USER=abc                   # Linuxserver user running amule (default: abc)
```

### Network Configuration

```bash
# P2P Network Ports
AMULE_PORT=4662                 # eMule TCP port (default: 4662)
AMULE_UDPPORT=4672              # eMule UDP port (default: 4672)
AMULE_UDP_ENABLE=1              # Enable UDP (1=true, 0=false, default: 1)
AMULE_YOURHOSTNAME=             # Your hostname (optional)
```

### Authentication & Security

```bash
# Password Configuration (choose one method)
AMULE_PASSWORD=mysecretpass     # Plain text password (will be MD5 hashed automatically)
AMULE_PASSWORD_HASH=6e0600b99a12d7e70687f0937a152fd5  # Pre-computed MD5 hash
```

### External Connect & Remote Control

```bash
AMULE_EC_ADDRESS=127.0.0.1      # External Connect IP (default: 127.0.0.1)
AMULE_EC_ECPORT=4712            # External Connect port (default: 4712)
```

### Web Server Configuration

```bash
AMULE_WEBSERVER_PORT=4711       # Web interface port (default: 4711)
AMULE_WEBSERVER_TEMPLATE=bootstrap  # Web template: bootstrap/reloaded/adaptable
```

### Behavior Control

```bash
# Configuration Generation
GENERATE_CONFIG=true            # Generate config (true/1/True/TRUE or unset=generate, false/0=don't)

# File Sharing Options  
RECURSIVE_SHARE=true            # Recursively share subdirectories (true/1 or false/0)

# Kad Network
FIX_KAD_GRAPH=true             # Fix Kad graph issues (true/1 or false/0)

# Auto-download Resources
DOWNLOAD_RESOURCES=true         # Auto-download nodes.dat, servers.met, ipfilter (true/1 or false/0)
```

## Default Values

The following default values are used when variables are not set:

- **AMULE_PORT**: 4662 (TCP P2P port)
- **AMULE_UDPPORT**: 4672 (UDP P2P port) 
- **AMULE_UDP_ENABLE**: 1 (UDP enabled)
- **AMULE_EC_ADDRESS**: 127.0.0.1 (External Connect address)
- **AMULE_EC_ECPORT**: 4712 (External Connect port)
- **AMULE_WEBSERVER_PORT**: 4711 (Web interface port)
- **AMULE_WEBSERVER_TEMPLATE**: bootstrap (Web UI template)
- **AMULE_PASSWORD_HASH**: 6e0600b99a12d7e70687f0937a152fd5 (MD5 hash for "amule")
- **GENERATE_CONFIG**: true (Generate configuration from template)
- **RECURSIVE_SHARE**: false (Don't share subdirectories recursively)
- **FIX_KAD_GRAPH**: false (Don't fix Kad graph)
- **DOWNLOAD_RESOURCES**: false (Don't auto-download resources)

## Configuration Files

The container automatically generates configuration files when `GENERATE_CONFIG=true` (default). You can also provide your own configuration files by mounting them to `/config`.

- **amule.conf**: Main configuration file
- **shareddir.dat**: Shared directories list (auto-generated if RECURSIVE_SHARE=true)
- **nodes.dat**: Kad nodes (auto-downloaded if DOWNLOAD_RESOURCES=true)
- **servers.met**: ED2K servers list (auto-downloaded if DOWNLOAD_RESOURCES=true)
- **ipfilter.dat**: IP filter (auto-downloaded if DOWNLOAD_RESOURCES=true)
 
## Web UI Templates

The container includes three different web interface templates:

- **bootstrap** (default): Modern responsive Bootstrap-based interface
- **reloaded**: Enhanced AmuleWebUI with additional features  
- **adaptable**: Adaptable web interface theme

Set via `AMULE_WEBSERVER_TEMPLATE` environment variable.

## Ports

The following ports are used by Amule:

| Port | Protocol | Description | Environment Variable |
|------|----------|-------------|---------------------|
| 4662 | TCP | eMule P2P network | AMULE_PORT |
| 4672 | UDP | eMule P2P network | AMULE_UDPPORT |
| 4711 | TCP | Web interface | AMULE_WEBSERVER_PORT |
| 4712 | TCP | External Connect (remote control) | AMULE_EC_ECPORT |

## Security Notes

- **Password Security**: Use `AMULE_PASSWORD` for automatic MD5 hashing, or provide your own `AMULE_PASSWORD_HASH`
- **External Access**: Change `AMULE_EC_ADDRESS` from `127.0.0.1` to `0.0.0.0` to allow external connections
- **Web Interface**: Accessible at `http://container-ip:4711` with configured password
- **Default Password**: The default password is "amule" (hash: 6e0600b99a12d7e70687f0937a152fd5)

## Volumes

- `/data`: Amule data directory (downloads, temp files)
- `/config`: Configuration files and runtime data

## Examples

### Basic Setup
```bash
docker run -d --name amule \
  -p 4662:4662 \
  -p 4672:4672/udp \
  -p 4711:4711 \
  -v ./amule-data:/data \
  -v ./amule-config:/config \
  -e AMULE_PASSWORD=mypassword \
  amule:latest
```

### Advanced Setup with Custom Configuration
```bash
docker run --rm -ti --name amule \
  -p 4662:4662 \
  -p 4672:4672/udp \
  -p 4711:4711 \
  -p 4712:4712 \
  -v ./data:/data \
  -v ./config:/config \
  -e AMULE_PASSWORD=supersecret \
  -e AMULE_WEBSERVER_TEMPLATE=reloaded \
  -e AMULE_EC_ADDRESS=0.0.0.0 \
  -e RECURSIVE_SHARE=true \
  -e DOWNLOAD_RESOURCES=true \
  -e FIX_KAD_GRAPH=false \
  amule:latest
```

# Author

Jose Riguera `<jriguera@gmail.com>`
