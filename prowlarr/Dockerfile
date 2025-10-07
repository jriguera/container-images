# syntax=docker/dockerfile:1

FROM ghcr.io/linuxserver/baseimage-alpine:3.22

# set version label
ARG BUILD_DATE="unknown"
ARG VERSION="latest"
ARG COMMIT_SHA="main"
ARG PROWLARR_VERSION
ARG PROWLARR_BRANCH="master"

LABEL org.opencontainers.image.created="${BUILD_DATE}"
LABEL org.opencontainers.image.authors="linuxserver.io authors and jriguera"
LABEL org.opencontainers.image.url="https://github.com/jriguera/container-images/tree/main/prowlarr"
LABEL org.opencontainers.image.documentation="https://github.com/jriguera/container-images/tree/main/prowlarr"
LABEL org.opencontainers.image.source="https://github.com/jriguera/container-images"
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.revision="${COMMIT_SHA}"
LABEL org.opencontainers.image.vendor="jriguera"
LABEL org.opencontainers.image.licenses="GPL-3.0-only"
LABEL org.opencontainers.image.ref.name="${COMMIT_SHA}"
LABEL org.opencontainers.image.title="Prowlarr"
LABEL org.opencontainers.image.description="[Prowlarr](https://github.com/Prowlarr/Prowlarr/) is an open-source file management tool."
LABEL build_version="Prowlarr version:- ${VERSION} Build-date:- ${BUILD_DATE}"
LABEL maintainer="jriguera"

ENV CONFIGDIR=/config
ENV PORT=9696
ENV COMPlus_EnableDiagnostics=0

# copy local files
COPY root/ /

RUN chmod a+x /healthcheck.sh && \
  echo "**** install runtime packages ****" && \
  apk add --no-cache \
    ca-certificates \
    icu-libs \
    sqlite-libs \
    xmlstarlet \
    jq \
    tzdata \
    socat \
    net-tools \
    util-linux \
    curl \
    bash \
  && \
  echo "**** cleanup ****" && \
  rm -rf /tmp/* $HOME/.cache

RUN \
  echo "**** downloading prowlarr ****" && \
  VERSION=${PROWLARR_VERSION:-$(curl -s https://api.github.com/repos/Prowlarr/Prowlarr/releases/latest | jq -r '.tag_name')} && PROWLARR_VERSION=${VERSION#v} && \
  curl -s -L "https://github.com/Prowlarr/Prowlarr/releases/download/v${PROWLARR_VERSION}/Prowlarr.${PROWLARR_BRANCH}.${PROWLARR_VERSION}.linux-musl-core-x64.tar.gz" -o /tmp/prowlarr.tar.gz && \
  mkdir /opt/prowlarr && tar -xzf /tmp/prowlarr.tar.gz -C /opt/prowlarr --strip-components=1 && \
  echo -e "UpdateMethod=docker\nBranch=${PROWLARR_BRANCH}\nPackageVersion=${VERSION}\nPackageAuthor=jriguera" > /opt/prowlarr/package_info && \
  printf "Prowlarr version: ${VERSION}\nBuild-date: ${BUILD_DATE}" > /build_version && \
  echo "**** cleanup ****" && \
  rm -rf /opt/prowlarr/Prowlarr.Update /tmp/*

# Define healthcheck script
HEALTHCHECK --start-period=2s --interval=60s --timeout=3s CMD /healthcheck.sh

# ports and volumes
EXPOSE "${PORT}"
VOLUME "${CONFIGDIR}"
