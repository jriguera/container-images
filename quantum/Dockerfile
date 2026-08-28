# syntax=docker/dockerfile:1

FROM ghcr.io/linuxserver/baseimage-alpine:3.24

# set version label
ARG BUILD_DATE="unknown"
ARG VERSION="latest"
ARG COMMIT_SHA="main"
ARG FILEBROWSER_VERSION="v2.0.2-beta"

LABEL org.opencontainers.image.created="${BUILD_DATE}"
LABEL org.opencontainers.image.authors="linuxserver.io authors and jriguera"
LABEL org.opencontainers.image.url="https://github.com/jriguera/container-images/tree/main/quantum"
LABEL org.opencontainers.image.documentation="https://github.com/jriguera/container-images/tree/main/quantum"
LABEL org.opencontainers.image.source="https://github.com/jriguera/container-images"
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.revision="${COMMIT_SHA}"
LABEL org.opencontainers.image.vendor="jriguera"
LABEL org.opencontainers.image.licenses="Apache-2.0"
LABEL org.opencontainers.image.ref.name="${COMMIT_SHA}"
LABEL org.opencontainers.image.title="FileBrowser Quantum"
LABEL org.opencontainers.image.description="[FileBrowser Quantum](https://github.com/gtsteffaniak/filebrowser) is a self-hosted web-based file manager (a feature-rich fork of File Browser)."
LABEL build_version="FileBrowser Quantum version:- ${VERSION} Build-date:- ${BUILD_DATE}"
LABEL maintainer="jriguera"

ENV DATADIR=/data
ENV CONFIGDIR=/config
ENV PORT=8000

# copy local files
COPY root/ /

RUN chmod a+x /healthcheck.sh && \
  echo "**** install runtime packages ****" && \
  apk add --no-cache \
    ca-certificates \
    mailcap \
    tzdata \
    socat \
    net-tools \
    curl \
    bash \
    ffmpeg \
  && \
  printf "FileBrowser Quantum version: ${VERSION}\nBuild-date: ${BUILD_DATE}" > /build_version && \
  echo "**** cleanup ****" && \
  rm -rf /tmp/* $HOME/.cache

# Download and install FileBrowser Quantum binary (single static Go binary)
RUN \
  echo "**** downloading filebrowser quantum ****" && \
  FB_VERSION=${FILEBROWSER_VERSION:-$(curl -s "https://api.github.com/repos/gtsteffaniak/filebrowser/releases/latest" | grep -oE '"tag_name": *"[^"]+"' | head -1 | cut -d'"' -f4)} && \
  case "$(uname -m)" in \
    x86_64)  FB_ARCH="amd64" ;; \
    aarch64) FB_ARCH="arm64" ;; \
    armv7l)  FB_ARCH="armv7" ;; \
    armv6l)  FB_ARCH="armv6" ;; \
    *) echo "Unsupported architecture: $(uname -m)" && exit 1 ;; \
  esac && \
  curl -s -L "https://github.com/gtsteffaniak/filebrowser/releases/download/${FB_VERSION}/linux-${FB_ARCH}-filebrowser" -o /bin/filebrowser && \
  chmod +x /bin/filebrowser && \
  echo "**** verifying binary ****" && \
  /bin/filebrowser version

# Define healthcheck script
HEALTHCHECK --start-period=5s --interval=60s --timeout=3s CMD /healthcheck.sh

# ports and volumes
EXPOSE "${PORT}"
VOLUME "${CONFIGDIR}" "${DATADIR}"
