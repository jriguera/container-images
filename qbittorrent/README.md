# docker-qbittorrent

qBittorrent Docker image based on Alpine, multi-arch.
qBittorrent is a free and open-source BitTorrent client that aims to provide a good alternative to µTorrent.

### Develop and test builds

Just type:

```
docker build \
    --no-cache \
    --pull \
    --build-arg BUILD_DATE=$(date '+%Y-%m-%dT%H:%M:%S%:z') \
    --build-arg VERSION=0.1 \
    --build-arg COMMIT_SHA=main \
    .  -t qbittorrent
```

# Usage

Given the docker image with name `qbittorrent`:

```bash
docker run --rm -ti --name qbittorrent \
  -p 6881:6881 \
  -p 6881:6881/udp \
  -p 8080:8080 \
  -v $(pwd)/data:/data \
  -v $(pwd)/config:/config \
  qbittorrent
```

Or using docker-compose:

```yaml
version: "3.9"
services:
  qbittorrent:
    image: qbittorrent:latest
    restart: unless-stopped
    ports:
      - "6881:6881"      # BitTorrent TCP
      - "6881:6881/udp"  # BitTorrent UDP  
      - "8080:8080"      # Web interface
    volumes:
      - ./data:/data
      - ./config:/config
    environment:
      - QBT_WEBUI_PASSWORD=mysecretpassword
      - QBT_WEBUI_PORT=8080
      - QBT_TORRENT_PORT=6881
```

You can access the web interface at `http://localhost:8080` using the configured password.

## Environment Variables

### Network Configuration

```bash
# BitTorrent Network
QBT_TORRENT_PORT=6881              # BitTorrent port (default: 6881)
QBT_TORRENT_IFACE_ADDR=0.0.0.0     # Interface address (default: 0.0.0.0)
```

### Web UI Configuration

```bash
QBT_WEBUI_PORT=8080                # Web interface port (default: 8080)
QBT_WEBUI_ENABLED=true             # Enable web interface (default: true)
QBT_WEBUI_PASSWORD=                # Web interface password (plain text)

# Themes
# Repository (format: "owner/repo"):  VueTorrent/VueTorrent or lantean-code/qbtmud
QBT_WEBUI_REPO=lantean-code/qbtmud
QBT_WEBUI_VERSION=                 # Version to install (leave empty for latest)
```

### Performance & Limits

```bash
# Memory and Connection Limits
QBT_MEMORY_WORKING_SET_LIMIT=1024  # Memory limit in MB (default: 1024)
QBT_GLOBAL_MAX_CONNECTIONS=200     # Max global connections (default: 200)
QBT_GLOBAL_MAX_UPLOADS=20          # Max global uploads (default: 20)

# Transfer Rate Limits
QBT_GLOBAL_MAX_DOWNLOAD_RATE=0     # Max download rate KB/s (0=unlimited, default: 0)
QBT_GLOBAL_MAX_UPLOAD_RATE=0       # Max upload rate KB/s (0=unlimited, default: 0)

# Active Torrent Limits
QBT_MAX_ACTIVE_DOWNLOADS=3         # Max active downloads (default: 3)
QBT_MAX_ACTIVE_UPLOADS=3           # Max active uploads (default: 3)
QBT_MAX_ACTIVE_TORRENTS=5          # Max active torrents (default: 5)
```

### Advanced Settings

```bash
# Configuration Generation
GENERATE_CONFIG=true               # Generate config (true/1/True/TRUE or unset=generate, false/0=don't)
```

## Default Values

The following default values are used when variables are not set:

- **QBT_WEBUI_PORT**: 8080 (Web interface port)
- **QBT_WEBUI_ENABLED**: true (Web interface enabled)
- **QBT_TORRENT_PORT**: 6881 (BitTorrent port)
- **QBT_TORRENT_IFACE_ADDR**: 0.0.0.0 (All interfaces)
- **QBT_MEMORY_WORKING_SET_LIMIT**: 1024 (MB)
- **QBT_GLOBAL_MAX_CONNECTIONS**: 200 (connections)
- **QBT_GLOBAL_MAX_UPLOADS**: 20 (uploads)
- **QBT_GLOBAL_MAX_DOWNLOAD_RATE**: 0 (unlimited KB/s)
- **QBT_GLOBAL_MAX_UPLOAD_RATE**: 0 (unlimited KB/s)
- **QBT_MAX_ACTIVE_DOWNLOADS**: 3 (torrents)
- **QBT_MAX_ACTIVE_UPLOADS**: 3 (torrents)
- **QBT_MAX_ACTIVE_TORRENTS**: 5 (torrents)
- **QBT_WEBUI_REPO**: lantean-code/qbtmud
- **QBT_WEBUI_VERSION**: latest
- **GENERATE_CONFIG**: true (Generate configuration from template)

## Configuration Files

The container automatically generates configuration files when `GENERATE_CONFIG=true` (default). You can also provide your own configuration files by mounting them to `/config`.

- **qBittorrent.conf**: Main configuration file
- **logs/**: Application logs directory
- **webui/**: Web interface files

## Password Management

qBittorrent uses PBKDF2 hashing for passwords. The container includes a `pbkdf2` utility for generating password hashes:

```bash
# Generate password hash
docker exec qbittorrent pbkdf2 "mypassword"
```

The `QBT_WEBUI_PASSWORD` environment variable should contain the plain text password, which will be automatically hashed during configuration generation.

## Ports

The following ports are used by qBittorrent:

| Port | Protocol | Description | Environment Variable |
|------|----------|-------------|---------------------|
| 6881 | TCP/UDP | BitTorrent network | QBT_TORRENT_PORT |
| 8080 | TCP | Web interface | QBT_WEBUI_PORT |

## Security Notes

- **Password Security**: Use `QBT_WEBUI_PASSWORD` for plain text passwords (automatically hashed)
- **Interface Binding**: Change `QBT_TORRENT_IFACE_ADDR` to bind to specific interfaces
- **Web Interface**: Accessible at `http://container-ip:8080` with configured password
- **Default Credentials**: Username: `admin`, Password: (set via `QBT_WEBUI_PASSWORD`)

## Volumes

- `/data`: qBittorrent data directory (downloads, temp files, torrents)
- `/config`: Configuration files and runtime data

## Examples

### Basic Setup
```bash
docker run -d --name qbittorrent \
  -p 6881:6881 \
  -p 6881:6881/udp \
  -p 8080:8080 \
  -v ./qbt-data:/data \
  -v ./qbt-config:/config \
  -e QBT_WEBUI_PASSWORD=mypassword \
  qbittorrent:latest
```

### Advanced Setup with Custom Configuration
```bash
docker run -d --name qbittorrent \
  -p 6881:6881 \
  -p 6881:6881/udp \
  -p 8080:8080 \
  -v ./data:/data \
  -v ./config:/config \
  -e QBT_WEBUI_PASSWORD=supersecret \
  -e QBT_WEBUI_PORT=8080 \
  -e QBT_TORRENT_PORT=6881 \
  -e QBT_GLOBAL_MAX_CONNECTIONS=500 \
  -e QBT_MAX_ACTIVE_DOWNLOADS=5 \
  -e QBT_MAX_ACTIVE_UPLOADS=3 \
  -e QBT_MEMORY_WORKING_SET_LIMIT=2048 \
  -e QBT_WEBUI_REPO=VueTorrent/VueTorrent \
  qbittorrent:latest
```
