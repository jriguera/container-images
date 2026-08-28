# syntax=docker/dockerfile:1

# Container versions
ARG NODE_VERSION=24.16.0
ARG COMPOSER_VERSION=2.10.1
ARG TA_VERSION=v5.0.0

FROM docker.io/library/node:${NODE_VERSION}-alpine AS node-builder

ARG TA_VERSION
RUN apk add --no-cache curl jq
RUN echo "**** downloading tasmoadmin ****" && \
  TA_VERSION=${TA_VERSION:-$(curl -s "https://api.github.com/repos/TasmoAdmin/TasmoAdmin/releases/latest" | jq -r '.tag_name')} && \
  curl -fsSL "https://github.com/TasmoAdmin/TasmoAdmin/releases/download/${TA_VERSION}/tasmoadmin_${TA_VERSION}.tar.gz" -o /tmp/tasmoadmin.tar.gz && \
  mkdir -p /tasmoadmin && tar -xzf /tmp/tasmoadmin.tar.gz --strip-components=1 -C /tasmoadmin && \
  rm -rf /tmp/tasmoadmin.tar.gz
WORKDIR /tasmoadmin
RUN echo "**** build ****" && \
  npm ci && npm run build && \
  touch .dockerenv && \
  rm -rf node_modules/


FROM docker.io/library/composer:${COMPOSER_VERSION} AS builder

COPY --from=node-builder /tasmoadmin  /tasmoadmin
WORKDIR /tasmoadmin
RUN composer install --optimize-autoloader --no-dev


FROM ghcr.io/linuxserver/baseimage-alpine:3.24

# set version label
ARG BUILD_DATE="unknown"
ARG VERSION="latest"
ARG COMMIT_SHA="main"
ARG TA_VERSION

LABEL org.opencontainers.image.created="${BUILD_DATE}"
LABEL org.opencontainers.image.authors="linuxserver.io authors and jriguera"
LABEL org.opencontainers.image.url="https://github.com/jriguera/container-images/tree/main/tasmoadmin"
LABEL org.opencontainers.image.documentation="https://github.com/jriguera/container-images/tree/main/tasmoadmin"
LABEL org.opencontainers.image.source="https://github.com/jriguera/container-images"
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.revision="${COMMIT_SHA}"
LABEL org.opencontainers.image.vendor="jriguera"
LABEL org.opencontainers.image.licenses="GPL-3.0-only"
LABEL org.opencontainers.image.ref.name="${COMMIT_SHA}"
LABEL org.opencontainers.image.title="TasmoAdmin"
LABEL org.opencontainers.image.description="[TasmoAdmin](https://github.com/TasmoAdmin/TasmoAdmin) is an open-source web interface for managing and monitoring Tasmota devices."
LABEL build_version="TasmoAdmin version:- ${VERSION} Build-date:- ${BUILD_DATE}"
LABEL maintainer="jriguera"

# environment variables
ENV CONFIGDIR=/config
ENV PORT=8088
ENV BUILD_VERSION=${VERSION}
ENV BUILD_REF=${COMMIT_SHA}
ENV SSL=false
ENV PHP_INI_SCAN_DIR=":/config/php"
ENV NGINX_USER=abc

RUN echo "**** install runtime packages ****" && \
  apk add --no-cache \
    libcrypto3 \
    libssl3 \
    busybox \
    musl-utils \
    musl \
    nginx \
        php85-curl \
        php85-fpm \
        php85-dom \
        php85-json \
        php85-openssl \
        php85-session \
        php85-zip \
        php85-mbstring \
        php85 \
    apache2-utils \
    logrotate \
    openssl \
    ca-certificates \
    pwgen \
    tzdata \
    socat \
    net-tools \
    curl \
    bash \
    gettext \
    jq \
  && \
  printf "Linuxserver.io version: ${VERSION}\nBuild-date: ${BUILD_DATE}" > /build_version && \
  echo "**** cleanup ****" && \
  rm -rf /tmp/* $HOME/.cache

# copy local files
COPY root/ /

# Setup application
COPY --from=builder /tasmoadmin /var/www/tasmoadmin

RUN echo "**** setup healthcheck ****" && \
  chmod a+x /healthcheck.sh && \
  echo "**** configure nginx ****" && \
  rm -f /etc/nginx/http.d/* && \
  echo "include ${CONFIGDIR}/nginx/nginx.conf;" > /etc/nginx/nginx.conf && \
  echo "**** configure php ****" && \
  ln -sf /usr/bin/php85 /usr/bin/php && \
  echo "**** fix logrotate ****" && \
  chmod a+x /etc/periodic/daily/logrotate && \
  echo "**** setup tasmoadmin symlink ****" && \
  rm -rf /var/www/tasmoadmin/data && ln -sf ${CONFIGDIR}/tasmoadmin /var/www/tasmoadmin/data && \
  touch /var/www/tasmoadmin/.dockerenv

# Define healthcheck script
HEALTHCHECK --start-period=2s --interval=60s --timeout=3s CMD /healthcheck.sh

# ports and volumes
EXPOSE "${PORT}"
VOLUME "${CONFIGDIR}"
