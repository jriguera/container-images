# docker-homepage

Homepage Dashboard Docker image based on Alpine with LinuxServer.io base, multi-arch.

[Homepage](https://gethomepage.dev/) is a modern, fully static, fast, secure, and highly customizable application dashboard with integrations for over 100 services and translations into multiple languages. Easily configured via YAML files or through docker labels.

## Features

- **Customizable Dashboard**: Organize your services with bookmarks, widgets, and integrations
- **Service Widgets**: Display real-time stats from 100+ services
- **Docker Integration**: Auto-discover and display running containers
- **Multi-language Support**: Available in many languages
- **Custom Patches**: This image includes a custom aMule widget for P2P monitoring

## Quick Start

### Using Docker CLI

```bash
docker run -d \
  --name=homepage \
  -e PUID=1000 \
  -e PGID=1000 \
  -e TZ=Etc/UTC \
  -p 3000:3000 \
  -v $(pwd)/config:/config \
  --restart unless-stopped \
  homepage
```

### Using Docker Compose

See the `docker-compose.yml` file in this directory for a complete example.

Access the dashboard at `http://localhost:3000`

## Build

To build the image locally:

```bash
docker build \
    --no-cache \
    --pull \
    --build-arg BUILD_DATE=$(date '+%Y-%m-%dT%H:%M:%S%:z') \
    --build-arg VERSION=0.1 \
    --build-arg COMMIT_BRANCH=main \
    . -t homepage
```

### Build Arguments

| Argument | Default | Description |
|----------|---------|-------------|
| `BUILD_DATE` | `unknown` | Build timestamp |
| `VERSION` | `latest` | Version tag for the image |
| `COMMIT_BRANCH` | `main` | Homepage git branch/tag to build from |

## Environment Variables

### Core Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PUID` | `911` | User ID for running the application |
| `PGID` | `911` | Group ID for running the application |
| `TZ` | `Etc/UTC` | Timezone (e.g., `America/New_York`) |

## Configuration

Homepage uses YAML configuration files stored in the `/config` directory:

### Configuration Files

| File | Purpose |
|------|---------|
| `settings.yaml` | Main settings (title, theme, language, etc.) |
| `services.yaml` | Service definitions and widgets |
| `bookmarks.yaml` | Bookmark definitions |
| `widgets.yaml` | Dashboard widgets (weather, search, etc.) |
| `docker.yaml` | Docker integration settings |
| `kubernetes.yaml` | Kubernetes integration settings |

### Auto-Generated Configuration

On first run, the container will generate default configuration files from templates located in `/defaults/`:
- `settings.json.template` → `/config/settings.json`

You can customize these files directly or mount your own configuration files.

### Configuration via Environment Variables

Some settings can be controlled via environment variables. See the [Homepage documentation](https://gethomepage.dev/configs/) for details.

## Volumes

| Path | Description |
|------|-------------|
| `/config` | Homepage configuration files (YAML) |

## Ports

| Port | Protocol | Description |
|------|----------|-------------|
| 3000 | TCP | Web interface |

## Custom Features

### aMule Widget

This image includes a custom aMule widget for monitoring aMule P2P client stats. To use it, add the following to your `services.yaml`:

```yaml
- aMule:
    icon: a-mule.png
    href: http://amule-host:4711
    description: aMule P2P Client
    widget:
      type: amule
      server: amule-host:4712
      password: your-ec-password
```

The widget displays:
- Download speed
- Upload speed
- Number of files downloading
- Queue size

See `/patches/amule.patch` for implementation details.

## Docker Integration

Homepage can auto-discover and display your running Docker containers. To enable this:

1. Mount the Docker socket:
   ```yaml
   volumes:
     - /var/run/docker.sock:/var/run/docker.sock:ro
   ```

2. Configure in `docker.yaml`:
   ```yaml
   my-docker:
     host: /var/run/docker.sock
   ```

3. Add container labels for automatic widget discovery:
   ```yaml
   labels:
     - homepage.group=Media
     - homepage.name=Plex
     - homepage.icon=plex.png
     - homepage.href=http://plex:32400
     - homepage.description=Media Server
   ```

## Examples

### Minimal Setup

```bash
docker run -d \
  --name=homepage \
  -p 3000:3000 \
  -v $(pwd)/config:/config \
  homepage
```

### With Docker Integration

```bash
docker run -d \
  --name=homepage \
  -e PUID=1000 \
  -e PGID=1000 \
  -e TZ=America/New_York \
  -p 3000:3000 \
  -v $(pwd)/config:/config \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  homepage
```

### Using Docker Compose

```yaml
version: "3.9"
services:
  homepage:
    image: homepage:latest
    container_name: homepage
    restart: unless-stopped
    ports:
      - "3000:3000"
    volumes:
      - ./config:/config
      - /var/run/docker.sock:/var/run/docker.sock:ro
    environment:
      - PUID=1000
      - PGID=1000
      - TZ=Etc/UTC
```

## Links

- [Homepage Official Documentation](https://gethomepage.dev/)
- [Homepage GitHub](https://github.com/gethomepage/homepage)
- [LinuxServer.io Base Images](https://github.com/linuxserver/docker-baseimage-alpine)
