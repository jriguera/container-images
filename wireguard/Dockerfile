# syntax=docker/dockerfile:1

# Container versions
ARG GOLANG_VERSION=1.25.5
FROM docker.io/library/golang:${GOLANG_VERSION}-alpine AS go-builder
ARG VERSION="latest"
ARG COMMIT_SHA="main"
RUN apk add --no-cache make git ca-certificates
# Download Go modules
COPY container-network /app
# Set destination for COPY
WORKDIR /app
RUN go mod download
RUN make build VERSION="${VERSION}"


FROM docker.io/serjs/go-socks5-proxy:latest AS socks5proxy


FROM ghcr.io/linuxserver/baseimage-alpine:3.22
# set version label
ARG BUILD_DATE="unknown"
ARG VERSION="latest"
ARG COMMIT_SHA="main"
ARG WIREGUARD_VERSION=1.0

LABEL org.opencontainers.image.created="${BUILD_DATE}"
LABEL org.opencontainers.image.authors="linuxserver.io authors and jriguera"
LABEL org.opencontainers.image.url="https://github.com/jriguera/container-images/tree/main/wireguard"
LABEL org.opencontainers.image.documentation="https://github.com/jriguera/container-images/tree/main/wireguard"
LABEL org.opencontainers.image.source="https://github.com/jriguera/container-images"
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.revision="${COMMIT_SHA}"
LABEL org.opencontainers.image.vendor="jriguera"
LABEL org.opencontainers.image.licenses="GPL-3.0-only"
LABEL org.opencontainers.image.ref.name="${COMMIT_SHA}"
LABEL org.opencontainers.image.title="wireguard"
LABEL org.opencontainers.image.description="[WireGuard](https://www.wireguard.com/) is an open-source VPN that utilizes state-of-the-art cryptography."
LABEL build_version="WireGuard version:- ${VERSION} Build-date:- ${BUILD_DATE}"
LABEL maintainer="jriguera"

# environment variables
ENV CONFIGDIR=/config
ENV PORT=51820
ENV PROXY_PORT=1080
ENV WEBHOOK_PORT=9000

# copy local files
COPY root/ /

RUN chmod a+x /healthcheck.sh && \
  echo "**** install dependencies ****" && \
  apk add --no-cache \
    bc \
    coredns \
    grep \
    iproute2 \
    iptables \
    ip6tables \
    iputils \
    tzdata \
    socat \
    kmod \
    libcap-utils \
    libqrencode-tools \
    net-tools \
    nftables \
    openresolv \
    curl \
    bash \
    jq \
    wireguard-tools~${WIREGUARD_VERSION} \
    && \
    echo "**** setting permissions on /usr/local/bin ****" && \
    chmod +x /usr/local/bin/*.sh && \
    echo "**** setting up wireguard config symlink ****" && \
    echo "wireguard" >> /etc/modules && \
    sed -i 's|\[\[ $proto == -4 \]\] && cmd sysctl -q net\.ipv4\.conf\.all\.src_valid_mark=1|[[ $proto == -4 ]] \&\& [[ $(sysctl -n net.ipv4.conf.all.src_valid_mark) != 1 ]] \&\& cmd sysctl -q net.ipv4.conf.all.src_valid_mark=1|' /usr/bin/wg-quick && \
    rm -rf /etc/wireguard && ln -s /config/wg_confs /etc/wireguard && \
    printf "Linuxserver.io version: ${VERSION}\nBuild-date: ${BUILD_DATE}" > /build_version && \
    echo "**** cleanup ****" && \
    rm -rf /tmp/* $HOME/.cache

# add socks5 proxy
COPY --from=socks5proxy /socks5 /usr/bin/socks5
COPY --from=go-builder /app/build/container-network /usr/bin/container-network

# Download and install webhook binary
RUN \
  echo "**** downloading webhook ****" && \
  WEBHOOK_VERSION=${WEBHOOK_VERSION:-$(curl -s "https://api.github.com/repos/adnanh/webhook/releases/latest" | jq -r '.tag_name')} && \
  curl -s -L "https://github.com/adnanh/webhook/releases/download/${WEBHOOK_VERSION}/webhook-linux-amd64.tar.gz" -o /tmp/webhook.tar.gz && \
  mkdir /tmp/webhook && tar -xzf /tmp/webhook.tar.gz -C /tmp/webhook --strip-components=1 && \
  mv /tmp/webhook/webhook /usr/bin/webhook && chmod +x /usr/bin/webhook && \
  rm -rf /tmp/webhook.tar.gz /tmp/webhook

# Define healthcheck script
HEALTHCHECK --start-period=30s --interval=60s --timeout=5s CMD /healthcheck.sh

# ports and volumes
EXPOSE "${PORT}/udp" "${PROXY_PORT}"
VOLUME "${CONFIGDIR}"
