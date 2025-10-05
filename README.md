# Container Images

This repository contains Docker container images for various applications, all based on the [LinuxServer.io](https://www.linuxserver.io/) base images. Each container inherits the standard LinuxServer.io environment variables and configuration patterns.

## Available Containers

| Container | Description | README |
|-----------|-------------|---------|
| [amule](amule) | aMule eDonkey/Kademlia P2P client | [README](amule/README.md) |
| [filebrowser](filebrowser) | Web-based file manager | [README](filebrowser/README.md) |
| [mosquitto](mosquitto) | Eclipse Mosquitto MQTT broker | [README](mosquitto/README.md) |
| [qbittorrent](qbittorrent) | qBittorrent BitTorrent client | [README](qbittorrent/README.md) |

## Common Environment Variables

Since all containers are based on LinuxServer.io base images, they share the following common environment variables:

### User/Group Identifiers
- `PUID=1000` - User ID for the container user
- `PGID=1000` - Group ID for the container user

### Timezone
- `TZ=Etc/UTC` - Specify a timezone to use (e.g., `Europe/London`, `America/New_York`)

### Additional Variables
Some containers may also support:
- `UMASK=022` - Control permissions of files and directories created

## Usage

Each container includes:
- A `docker-compose.yml` file for easy deployment
- A `Dockerfile` for building the image
- Configuration templates and defaults
- Health check scripts
- S6-overlay service definitions

## Getting Started

1. Navigate to the desired container directory
2. Review the container-specific README for configuration details
3. Modify the `docker-compose.yml` file as needed
4. Run with `docker-compose up -d`

## Base Image Features

All containers inherit LinuxServer.io base image features:
- S6-overlay for process supervision
- Automatic user/group creation
- Configuration file templating
- Health checks

For more information about LinuxServer.io base images, visit: https://docs.linuxserver.io/

# Release Process

The workflow will build, push with latest tag, and create a release when:

* Make changes to the image container directory
* Commit and push to trigger a build
* Create an annotated tag eg: `git tag -a amule-v1.0.0 -m "aMule v1.0.0"`
* Push the tag, eg: `git push origin amule-v1.0.0`
