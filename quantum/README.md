# docker-quantum

FileBrowser Quantum Docker image based on Alpine.
[FileBrowser Quantum](https://github.com/gtsteffaniak/filebrowser) is a
self-hosted web-based file manager. It is a feature-rich fork of the original
[File Browser](https://github.com/filebrowser/filebrowser) with multiple
sources, a modern responsive UI, efficient indexed search, OIDC/LDAP/proxy
auth, WebDAV and rich media previews.

### Develop and test builds

Just type:

```
docker build \
    --no-cache \
    --pull \
    --build-arg BUILD_DATE=$(date '+%Y-%m-%dT%H:%M:%S%:z') \
    --build-arg VERSION=0.1 \
    --build-arg COMMIT_SHA=main \
    .  -t quantum
```

To build against a different release (for example the stable v1.5 line):

```
docker build --build-arg FILEBROWSER_VERSION=v1.5.3-stable . -t quantum
```

# Usage

Given the docker image with name `quantum`:

```
docker run --rm -ti --name quantum \
    -p 8000:8000 \
    -e FB_ADMIN_USER="admin" \
    -e FB_ADMIN_PASSWORD="admin" \
    -v $(pwd)/data:/data \
    -v $(pwd)/config:/config \
    quantum
```

Then open http://localhost:8000 and log in with the admin credentials above.

You can also use environment variables to automatically define some settings:

```
PORT=8000
FB_BASE_URL="/"
FB_SOURCE="/data"
FB_ADMIN_USER="admin"
FB_ADMIN_PASSWORD="admin"
FB_LOG_LEVEL="info"
RESET_DB="false"
```

And use them:

```
docker run --name quantum -p 8000:8000 \
    -v $(pwd)/data:/data -v $(pwd)/config:/config \
    -e FB_ADMIN_USER="admin" -e FB_ADMIN_PASSWORD="secret" \
    -e FB_BASE_URL="/files" \
    -d ghcr.io/jriguera/container-images/quantum:latest
```

## Variables

On startup the container renders `/config/config.yaml` from the environment
variables below (unless `FB_CONFIG_FILE` points to your own file). See the
init script and [the official docs](https://filebrowserquantum.com/en/docs/configuration/)
for more details.

```
# FileBrowser Quantum configuration parameters and defaults
PORT=8000                                   # container listening port
FB_PORT=${FB_PORT:-${PORT}}
FB_BASE_URL="${FB_BASE_URL:-/}"             # subpath when behind a reverse proxy
FB_SOURCE="${FB_SOURCE:-/data}"             # root folder exposed as a source
FB_DATABASE="${FB_DATABASE:-/config/filebrowser.sqlite}"
FB_CACHE_DIR="${FB_CACHE_DIR:-/config/tmp}" # thumbnails / temp files
FB_CONFIG_FILE="${FB_CONFIG_FILE:-/config/config.yaml}"
FB_LOG_LEVEL="${FB_LOG_LEVEL:-info}"
FB_ADMIN_USER="${FB_ADMIN_USER:-admin}"
FB_ADMIN_PASSWORD="${FB_ADMIN_PASSWORD:-admin}"
FB_PASSWORD_MIN_LENGTH="${FB_PASSWORD_MIN_LENGTH:-2}"
FB_FFMPEG_PATH="${FB_FFMPEG_PATH:-/usr/bin}"  # ffmpeg/ffprobe dir for media thumbnails ("" disables)
RESET_DB="${RESET_DB:-false}"              # remove the database on startup
```

## Using your own configuration

Mount a YAML file and point `FB_CONFIG_FILE` at it to bypass the generated
configuration entirely:

```
docker run --name quantum -p 8000:8000 \
    -v $(pwd)/data:/data -v $(pwd)/config:/config \
    -v $(pwd)/config.yaml:/config/myconfig.yaml \
    -e FB_CONFIG_FILE="/config/myconfig.yaml" \
    quantum
```

Secrets like admin password can also be provided via
FileBrowser Quantum's own environment variables such as
`FILEBROWSER_ADMIN_PASSWORD`. See the
[environment variables reference](https://filebrowserquantum.com/en/docs/reference/environment-variables/).

## Notes

- The image defaults to the latest **v2 beta** Quantum release
  (`v2.0.2-beta`). To use the stable v1.5 line instead, build with
  `--build-arg FILEBROWSER_VERSION=v1.5.3-stable` (note the config schema
  differs: v1.5 keeps `port`/`baseURL` under `server:`).
- HTTP options (`port`, `baseURL`, ...) live under the top-level `http:` key in
  the v2 config; server options (`database.path`, `cacheDir`, `logging`,
  `sources`) live under `server:`. The database is SQLite
  (`filebrowser.sqlite`).
- The healthcheck queries the `/health` endpoint on the configured port and
  base URL.
- `ffmpeg`/`ffprobe` are bundled and enabled via `integrations.media.ffmpegPath`
  (`/usr/bin`) so media, video and album-art thumbnails work out of the box.
  Set `FB_FFMPEG_PATH=""` to disable the media integration.
