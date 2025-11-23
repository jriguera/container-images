# syntax=docker/dockerfile:1

# Container versions
ARG NODE_VERSION=22.20.0
ARG COMPOSER_VERSION=2.8.12

FROM docker.io/library/node:${NODE_VERSION}-alpine AS node-builder

RUN apk add --no-cache curl jq
RUN echo "**** downloading tasmoadmin ****" && \
  TA_VERSION=${TA_VERSION:-$(curl -s "https://api.github.com/repos/TasmoAdmin/TasmoAdmin/releases/latest" | jq -r '.tag_name')} && \
  curl -s -L "https://github.com/TasmoAdmin/TasmoAdmin/releases/download/${TA_VERSION}/tasmoadmin_${TA_VERSION}.tar.gz" -o /tmp/tasmoadmin.tar.gz && \
  mkdir -p /tasmoadmin && tar -xzf /tmp/tasmoadmin.tar.gz --strip-components=1 -C /tasmoadmin && \
  rm -rf /tmp/tasmoadmin.tar.gz
WORKDIR /tasmoadmin
RUN echo "**** build ****" && \
  npm ci && npm run build && \
  touch .dockerenv && \
  rm -rf node_modules/


# syntax=docker/dockerfile:1
FROM docker.io/library/composer:${COMPOSER_VERSION} AS builder

COPY --from=node-builder /tasmoadmin  /tasmoadmin
WORKDIR /tasmoadmin
RUN composer install --optimize-autoloader --no-dev


# syntax=docker/dockerfile:1
FROM ghcr.io/linuxserver/baseimage-alpine:3.22

# set version label
ARG BUILD_DATE="unknown"
ARG VERSION="latest"
ARG COMMIT_SHA="main"
ARG TA_VERSION=4.3.2

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
    php84-curl \
    php84-fpm \
    php84-dom \
    php84-json \
    php84-session \
    php84-zip \
    php84-mbstring \
    php84 \
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
  [ "$(readlink /usr/bin/php)" != "php84" ] && rm -rf /usr/bin/php && ln -s /usr/bin/php84 /usr/bin/php && \
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
