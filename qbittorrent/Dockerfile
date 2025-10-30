# syntax=docker/dockerfile:1

FROM ghcr.io/linuxserver/unrar:latest AS unrar

FROM ghcr.io/linuxserver/baseimage-alpine:edge

# set version label
ARG BUILD_DATE="unknown"
ARG VERSION="latest"
ARG COMMIT_SHA="main"
ARG QBITTORRENT_VERSION="5.1"
ARG QBITTORRENT_CLI_VERSION="2.2"

LABEL org.opencontainers.image.created="${BUILD_DATE}"
LABEL org.opencontainers.image.authors="linuxserver.io authors and jriguera"
LABEL org.opencontainers.image.url="https://github.com/jriguera/container-images/tree/main/qbittorrent"
LABEL org.opencontainers.image.documentation="https://github.com/jriguera/container-images/tree/main/qbittorrent"
LABEL org.opencontainers.image.source="https://github.com/jriguera/container-images"
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.revision="${COMMIT_SHA}"
LABEL org.opencontainers.image.vendor="jriguera"
LABEL org.opencontainers.image.licenses="GPL-3.0-only"
LABEL org.opencontainers.image.ref.name="${COMMIT_SHA}"
LABEL org.opencontainers.image.title="QBittorrent"
LABEL org.opencontainers.image.description="[QBittorrent](https://www.qbittorrent.org/) is a free P2P file sharing application."
LABEL build_version="QBittorrent version:- ${VERSION} Build-date:- ${BUILD_DATE}"
LABEL maintainer="jriguera"

# environment variables
ENV DATADIR=/downloads
ENV CONFIGDIR=/config
ENV PORT=8080

# add unrar
COPY --from=unrar /usr/bin/unrar-alpine /usr/bin/unrar

# copy local files
COPY root/ /

RUN chmod a+x /healthcheck.sh && \
  echo "**** install runtime packages ****" && \
  apk add --no-cache \
    qbittorrent-nox~${QBITTORRENT_VERSION} \
    grep \
    icu-libs \
    p7zip \
    python3 \
    qt6-qtbase-sqlite \
    ca-certificates \
    pwgen \
    su-exec \
    tzdata \
    socat \
    net-tools \
    curl \
    bash \
  && \
  apk add --no-cache --repository http://dl-cdn.alpinelinux.org/alpine/edge/testing \
    qbittorrent-cli~${QBITTORRENT_CLI_VERSION} \
  && \
  printf "Linuxserver.io version: ${VERSION}\nBuild-date: ${BUILD_DATE}" > /build_version && \
  echo "**** cleanup ****" && \
  rm -rf /tmp/* $HOME/.cache


# Define healthcheck script
HEALTHCHECK --start-period=2s --interval=60s --timeout=3s CMD /healthcheck.sh

# ports and volumes
EXPOSE "${PORT}"
VOLUME "${DATADIR}" "${CONFIGDIR}"

# symlink config to home to avoid additional directory level in /config
RUN ln -s ${CONFIGDIR} qBittorrent
