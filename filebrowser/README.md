# docker-filebrowser

Filebrowser Docker image based on Alpine, multi-arch.
[Filebrowser](https://github.com/filebrowser/filebrowser/) is an open-source file management tool.

### Develop and test builds

Just type:

```
docker build \
    --no-cache \
    --pull \
    --build-arg BUILD_DATE=$(date '+%Y-%m-%dT%H:%M:%S%:z') \
    --build-arg VERSION=0.1 \
    --build-arg COMMIT_SHA=main \
    .  -t filebrowser
```

# Usage

Given the docker image with name `filebrowser` from Github Package Repository:

```
docker run --rm -ti --name fb -p 8080:8080 -e FB_AUTH_METHOD="noauth" -e FB_USERS="admin:admin1" -e FB_PORT=8080 -v $(pwd)/data:/data -v $(pwd)/config:/config filebrowser
```

You can also use environment variables to automatically define some settings:

```
FB_PORT=8080
FB_BASE_URL=""
FB_ADDRESS=""
FB_LOG_DST="stdout"
FB_AUTH_METHOD="noauth"
RESET_DB="false"
FB_USERS="admin:a1 user2:b2:admin,execute"
```

And use them:

```
docker run --name fb -p 8080:8080 -v $(pwd)/data:/data -v $(pwd)/config:/config -e FB_AUTH_METHOD="noauth" -e FB_USERS="admin:admin1" -e FB_PORT=8080 -e FB_BASE_URL="/files" -d ghcr.io/jriguera/docker-filebrowser/filebrowser:latest
```

## Variables

See the configuration script for more details.

```
# Filebrowser configuration parameters and defaults
FB_PORT=${FB_PORT:-${PORT}}
FB_LOG_DST="${FB_LOG_DST:-stdout}"
FB_BASE_URL="${FB_BASE_URL:-}"
FB_ADDRESS="${FB_ADDRESS:-}"
FB_CONFIG_FILE="${FB_CONFIG_FILE:-/config/settings.json}"
```
