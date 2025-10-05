# syntax=docker/dockerfile:1

FROM ghcr.io/linuxserver/baseimage-alpine:3.22

# set version label
ARG BUILD_DATE="unknown"
ARG VERSION="latest"
ARG COMMIT_SHA="main"
ARG MOSQUITTO_VERSION=2

LABEL org.opencontainers.image.created="${BUILD_DATE}"
LABEL org.opencontainers.image.authors="linuxserver.io authors and jriguera"
LABEL org.opencontainers.image.url="https://github.com/jriguera/container-images/tree/main/mosquitto"
LABEL org.opencontainers.image.documentation="https://github.com/jriguera/container-images/tree/main/mosquitto"
LABEL org.opencontainers.image.source="https://github.com/jriguera/container-images"
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.revision="${COMMIT_SHA}"
LABEL org.opencontainers.image.vendor="jriguera"
LABEL org.opencontainers.image.licenses="GPL-3.0-only"
LABEL org.opencontainers.image.ref.name="${COMMIT_SHA}"
LABEL org.opencontainers.image.title="Eclipse Mosquitto"
LABEL org.opencontainers.image.description="[Eclipse Mosquitto](https://mosquitto.org/) is an open-source message broker that implements the MQTT protocol."
LABEL build_version="Eclipse Mosquitto version:- ${VERSION} Build-date:- ${BUILD_DATE}"
LABEL maintainer="jriguera"

# environment variables
ENV DATADIR=/data
ENV CONFIGDIR=/config
ENV PORT=1883
ENV MOSQUITTO_USER=abc

# copy local files
COPY root/ /

RUN chmod a+x /healthcheck.sh && \
  echo "**** install runtime packages ****" && \
  apk add --no-cache \
    mosquitto~${MOSQUITTO_VERSION} \
    mosquitto-libs~${MOSQUITTO_VERSION} \
    mosquitto-clients~${MOSQUITTO_VERSION} \
    ca-certificates \
    pwgen \
    su-exec \
    tzdata \
    socat \
    net-tools \
    curl \
    bash \
  && \
  printf "Linuxserver.io version: ${VERSION}\nBuild-date: ${BUILD_DATE}" > /build_version && \
  echo "**** cleanup ****" && \
  rm -rf /tmp/* $HOME/.cache

# Define healthcheck script
HEALTHCHECK --start-period=2s --interval=60s --timeout=3s CMD /healthcheck.sh

# ports and volumes
EXPOSE "${PORT}"
VOLUME "${DATADIR}"
