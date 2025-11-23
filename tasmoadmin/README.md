# docker-tasmoadmin

TasmoAdmin Docker image based on Alpine, multi-arch.
TasmoAdmin is an open-source web interface for managing and monitoring Tasmota devices.

### Develop and test builds

Just type:

```bash
docker build \
    --no-cache \
    --pull \
    --build-arg BUILD_DATE=$(date '+%Y-%m-%dT%H:%M:%S%:z') \
    --build-arg VERSION=0.1 \
    --build-arg COMMIT_SHA=main \
    .  -t tasmoadmin
```

# Usage

Given the docker image with name `tasmoadmin`:

```bash
docker run --rm -ti --name tasmoadmin \
  -p 8088:8088 \
  -v $(pwd)/config:/config \
  tasmoadmin
```

Or using docker-compose:

```yaml
version: "3.9"
services:
  tasmoadmin:
    image: tasmoadmin:latest
    restart: unless-stopped
    ports:
      - "8088:8088"      # Web interface
    volumes:
      - ./config:/config
    environment:
      - PUID=1000
      - PGID=1000
      - TZ=Etc/UTC
      - TA_USER="admin"
      - TA_PASSWORD="admin"
```

You can access the web interface at `http://localhost:8088`.

## Environment Variables

### TasmoAdmin Application Settings

```bash
TA_OTA_IP=""                               # IP address for OTA updates (optional)
TA_USER="admin"                            # TasmoAdmin username (default: admin)
TA_PASSWORD="admin"                        # TasmoAdmin password (default: admin)
```

### User/Group Identifiers (LinuxServer.io base)

```bash
PUID=1000                                  # User ID for the container user
PGID=1000                                  # Group ID for the container user
```

### Timezone

```bash
TZ=Etc/UTC                                 # Specify a timezone to use (e.g., Europe/London)
```

### Nginx Configuration

```bash
NGINX_SERVER_NAME="tasmoadmin.local"       # Server name for nginx (default: tasmoadmin.local)
NGINX_PORT="${PORT:-8080}"                 # Nginx port (defaults to PORT or 8080)
NGINX_SSL=false                            # Enable/disable SSL (true/false, default: false)
NGINX_CONFIG_FILE="nginx.conf.template"    # Nginx configuration template file
NGINX_CONFIG_FILE_SSL="ssl.conf.template"  # SSL configuration template file
```

### Configuration Generation

```bash
GENERATE_CONFIG=true                       # Generate configuration files from templates (true/false, default: true)
```

## Volumes

- `/config` - Contains TasmoAdmin configuration, device data, logs, and optional custom nginx/php configs
  - `/config/tasmoadmin` - TasmoAdmin data directory (MyConfig.json, devices.csv, etc.)
  - `/config/nginx` - Custom nginx configuration files
  - `/config/nginx/certs` - SSL certificates
  - `/config/nginx/site-confs` - Nginx site configurations
  - `/config/templates` - Configuration templates for customization
  - `/config/log` - Application and web server logs
  - `/config/php` - Custom PHP configuration files (ini files)

## Ports

- `8088` - Web interface (configurable via PORT environment variable)

## SSL/TLS Support

To enable SSL:

1. Set `NGINX_SSL=true` environment variable
2. Place your SSL certificates in `/config/nginx/certs/`:
   - `cert.crt` - SSL certificate
   - `cert.key` - SSL private key
   
Self-signed certificates can be generated using the template in `/config/nginx/certs/cert.conf`.

## Custom Configuration

Configuration files are automatically generated from templates if `GENERATE_CONFIG=true` (default).

### Nginx Configuration
You can customize nginx by:
1. Modifying templates in `/config/templates/`:
   - `nginx.conf.template` - Main nginx configuration
   - `ssl.conf.template` - SSL configuration
   - `default.conf.template` - Site configuration
2. Placing custom configuration files in `/config/nginx/site-confs/`

### PHP Configuration
Add custom PHP ini files to `/config/php/` directory.

### Environment File
You can create an `/config/env` file to set environment variables that will be loaded automatically.

## TasmoAdmin Configuration

The TasmoAdmin configuration is stored in `/config/tasmoadmin/MyConfig.json`. You can set the initial admin credentials using:
- `TA_USER` - Admin username (default: admin)
- `TA_PASSWORD` - Admin password (default: admin)

For OTA firmware updates, you can set the `TA_OTA_IP` variable to specify the IP address where firmware files are hosted.

## Features

- Multi-platform support (amd64, arm64, armv7)
- S6-overlay for process supervision
- Automatic user/group creation
- Health checks included
- Log rotation configured
- Based on LinuxServer.io base image
 