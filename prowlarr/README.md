# docker-prowlarr

Prowlarr Docker image based on Alpine, multi-arch.
Prowlarr is an indexer manager/proxy built on the popular *arr .net/reactjs base stack to integrate with your various PVR apps.

For more information, visit: https://wiki.servarr.com/prowlarr

For additional environment variables and advanced configuration options, see: https://wiki.servarr.com/prowlarr/environment-variables

### Develop and test builds

Just type:

```
docker build \
    --no-cache \
    --pull \
    --build-arg BUILD_DATE=$(date '+%Y-%m-%dT%H:%M:%S%:z') \
    --build-arg VERSION=0.1 \
    --build-arg COMMIT_SHA=main \
    .  -t prowlarr
```

Build parameters:
```
# Branch Selection (build-time)
PROWLARR_BRANCH=master            # Prowlarr branch to install (master, develop, nightly)
PROWLARR_VERSION=                 # Specific version to install (leave empty for latest)
```

# Usage

Given the docker image with name `prowlarr`:

```bash
docker run --rm -ti --name prowlarr \
  -p 9696:9696 \
  -v $(pwd)/config:/config \
  prowlarr
```

Or using docker-compose:

```yaml
version: "3.9"
services:
  prowlarr:
    image: prowlarr:latest
    restart: unless-stopped
    ports:
      - "9696:9696"      # Web interface
    volumes:
      - ./config:/config
    environment:
      - PROWLARR_PORT=9696
      - PROWLARR_APIKEY=your-api-key-here
```

You can access the web interface at `http://localhost:9696` without authentication

## Environment Variables

For a complete list of environment variables and advanced configuration options, refer to the official documentation: https://wiki.servarr.com/prowlarr/environment-variables

### Network Configuration

```bash
# Web Interface
PROWLARR_PORT=9696                 # Web interface port (default: 9696)
PROWLARR_URLBASE=                  # URL base for reverse proxy setup (default: empty)
```

### Authentication Configuration

```bash
PROWLARR_AUTH=Forms                # Authentication method: None, Forms, Basic (default: Forms)
PROWLARR_AUTH_MODE=DisabledForLocalAddresses  # Authentication mode: Enabled or DisabledForLocalAddresses (default: DisabledForLocalAddresses)
PROWLARR_APIKEY=                   # API key (auto-generated if empty)
```

### Advanced Settings

```bash
# Configuration Generation
GENERATE_CONFIG=true               # Generate config (true/1/True/TRUE or unset=generate, false/0=don't)
```

> **Note**: See https://wiki.servarr.com/prowlarr/environment-variables for additional environment variables including database configuration, logging options, and other advanced settings.

## Default Values

The following default values are used when variables are not set:

- **PROWLARR_PORT**: 9696 (Web interface port)
- **PROWLARR_AUTH**: Forms (Authentication method)
- **PROWLARR_AUTH_MODE**: DisabledForLocalAddresses
- **PROWLARR_URLBASE**: (empty - no URL base)
- **PROWLARR_APIKEY**: (auto-generated UUID if not provided)
- **GENERATE_CONFIG**: true (Generate configuration from template)

## Configuration Files

The container automatically generates configuration files when `GENERATE_CONFIG=true` (default). You can also provide your own configuration files by mounting them to `/config`.

- **config.xml**: Main configuration file
- **logs/**: Application logs directory
- **prowlarr.db**: SQLite database
- **Definitions/**: Indexer definitions (auto-updated)

## Environment File

You can create a `/config/env` file to set environment variables that will be loaded at container startup. This is useful for setting sensitive values like API keys:

```bash
# /config/env
PROWLARR_APIKEY=your-secret-api-key-here
PROWLARR_AUTH=Basic
```

## Authentication Methods

- **None**: No authentication required
- **Forms**: Web form-based authentication (default)

## Ports

The following port is used by Prowlarr:

| Port | Protocol | Description | Environment Variable |
|------|----------|-------------|---------------------|
| 9696 | TCP | Web interface | PROWLARR_PORT |

## Security Notes

- **API Key**: Use `PROWLARR_APIKEY` to set a custom API key, or let the container generate one automatically
- **Authentication**: Configure appropriate authentication method for your environment
- **Reverse Proxy**: Use `PROWLARR_URLBASE` when running behind a reverse proxy
- **Default Access**: Web interface accessible at `http://container-ip:9696`

## Volumes

- `/config`: Configuration files and runtime data

## Examples

### Basic Setup
```bash
docker run -d --name prowlarr \
  -p 9696:9696 \
  -v ./prowlarr-config:/config \
  prowlarr:latest
```

### Advanced Setup with Custom Configuration
```bash
docker run -d --name prowlarr \
  -p 9696:9696 \
  -v ./config:/config \
  -e PROWLARR_PORT=9696 \
  -e PROWLARR_AUTH=Basic \
  -e PROWLARR_APIKEY=your-custom-api-key \
  -e PROWLARR_URLBASE=/prowlarr \
  prowlarr:latest
```

### Behind Reverse Proxy
```bash
docker run -d --name prowlarr \
  -p 9696:9696 \
  -v ./config:/config \
  -e PROWLARR_URLBASE=/prowlarr \
  -e PROWLARR_AUTH=Forms \
  prowlarr:latest
```

## Integration with *arr Apps

Prowlarr can automatically configure indexers for your other *arr applications. After setting up Prowlarr:

1. Add your indexers in Prowlarr
2. Configure your *arr apps (Sonarr, Radarr, etc.) in Prowlarr's "Apps" section
3. Prowlarr will automatically sync indexers to your configured applications

## Health Check

The container includes a health check that verifies Prowlarr is responding on the configured port. The health check runs every 60 seconds with a 3-second timeout.