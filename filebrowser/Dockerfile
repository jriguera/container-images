# syntax=docker/dockerfile:1

FROM ghcr.io/linuxserver/baseimage-alpine:3.22

# set version label
ARG BUILD_DATE="unknown"
ARG VERSION="latest"
ARG COMMIT_SHA="main"
ARG FILEBROWSER_VERSION

LABEL org.opencontainers.image.created="${BUILD_DATE}"
LABEL org.opencontainers.image.authors="linuxserver.io authors and jriguera"
LABEL org.opencontainers.image.url="https://github.com/jriguera/container-images/tree/main/filebrowser"
LABEL org.opencontainers.image.documentation="https://github.com/jriguera/container-images/tree/main/filebrowser"
LABEL org.opencontainers.image.source="https://github.com/jriguera/container-images"
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.revision="${COMMIT_SHA}"
LABEL org.opencontainers.image.vendor="jriguera"
LABEL org.opencontainers.image.licenses="GPL-3.0-only"
LABEL org.opencontainers.image.ref.name="${COMMIT_SHA}"
LABEL org.opencontainers.image.title="Filebrowser"
LABEL org.opencontainers.image.description="[Filebrowser](https://github.com/filebrowser/filebrowser/) is an open-source file management tool."
LABEL build_version="Filebrowser version:- ${VERSION} Build-date:- ${BUILD_DATE}"
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
    jq \
    libcap \
    tzdata \
    socat \
    net-tools \
    curl \
    bash \
  && \
  printf "Filebrowser version: ${VERSION}\nBuild-date: ${BUILD_DATE}" > /build_version && \
  echo "**** cleanup ****" && \
  rm -rf /tmp/* $HOME/.cache

# Download and install filebrowser binary
RUN \
  echo "**** downloading filebrowser ****" && \
  FB_VERSION=${FILEBROWSER_VERSION:-$(curl -s "https://api.github.com/repos/filebrowser/filebrowser/releases/latest" | jq -r '.tag_name')} && \
  curl -s -L "https://github.com/filebrowser/filebrowser/releases/download/${FB_VERSION}/linux-amd64-filebrowser.tar.gz" -o /tmp/filebrowser.tar.gz && \
  mkdir /tmp/filebrowser && tar -xzf /tmp/filebrowser.tar.gz -C /tmp/filebrowser && \
  mv /tmp/filebrowser/filebrowser /bin/filebrowser && chmod +x /bin/filebrowser && \
  rm -rf /tmp/filebrowser.tar.gz /tmp/filebrowser

# Define healthcheck script
HEALTHCHECK --start-period=2s --interval=60s --timeout=3s CMD /healthcheck.sh

# ports and volumes
EXPOSE "${PORT}"
VOLUME "${CONFIGDIR}" "${DATADIR}"
