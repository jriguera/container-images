# syntax=docker/dockerfile:1

FROM alpine:edge AS builder

# Install packages
RUN \
  echo "**** install runtime packages ****" && \
  apk add --no-cache git unzip wget xz bash ca-certificates curl

RUN \
  echo "**** Download webui themes (shallow clones) ****" && \
  mkdir -p /build &&  cd /tmp && \
  git clone --depth 1 --single-branch https://github.com/pedro77/amuleweb-bootstrap-template.git bootstrap && \
  mv bootstrap/dist /build/bootstrap && \
  git clone --depth 1 --single-branch https://github.com/MatteoRagni/AmuleWebUI-Reloaded reloaded && \
  rm -rf reloaded/.git* reloaded/doc-images && mv reloaded /build/reloaded && \
  git clone --depth 1 --single-branch https://github.com/esaracho/amuleweb-adaptable adaptable && \
  rm -rf adaptable/.git* && mv adaptable /build/adaptable


# syntax=docker/dockerfile:1
FROM ghcr.io/linuxserver/baseimage-alpine:3.22

# set version label
ARG BUILD_DATE="unknown"
ARG VERSION="latest"
ARG COMMIT_SHA="main"
ARG AMULE_VERSION=2.3.3

LABEL org.opencontainers.image.created="${BUILD_DATE}"
LABEL org.opencontainers.image.authors="linuxserver.io authors and jriguera"
LABEL org.opencontainers.image.url="https://github.com/jriguera/container-images/tree/main/amule"
LABEL org.opencontainers.image.documentation="https://github.com/jriguera/container-images/tree/main/amule"
LABEL org.opencontainers.image.source="https://github.com/jriguera/container-images"
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.revision="${COMMIT_SHA}"
LABEL org.opencontainers.image.vendor="jriguera"
LABEL org.opencontainers.image.licenses="GPL-3.0-only"
LABEL org.opencontainers.image.ref.name="${COMMIT_SHA}"
LABEL org.opencontainers.image.title="Amule"
LABEL org.opencontainers.image.description="[Amule](https://www.amule.org/) is a free P2P file sharing application."
LABEL build_version="Amule version:- ${VERSION} Build-date:- ${BUILD_DATE}"
LABEL maintainer="jriguera"

# environment variables
ENV DATADIR=/data
ENV CONFIGDIR=/config
ENV PORT=4711
ENV ECPORT=4712
ENV AMULE_USER=abc

# Install packages
RUN \
  echo "**** install runtime packages ****" && \
  apk add --no-cache \
    libedit \
    libgcc \
    libintl \
    libpng \
    libstdc++ \
    libupnp \
    musl \
    wxwidgets \
    zlib \
    tzdata \
    pwgen \
    mandoc \
    curl \
    coreutils \
    bash \
    ca-certificates \
    socat \
    net-tools \
  && \
  apk add --no-cache --repository=http://dl-cdn.alpinelinux.org/alpine/edge/testing \
    crypto++ \
    amule~${AMULE_VERSION} \
    amule-doc~${AMULE_VERSION} \
  && \
  printf "Amule version: ${VERSION}\nBuild-date: ${BUILD_DATE}" > /build_version && \
  echo "**** cleanup ****" && \
  rm -rf /tmp/* $HOME/.cache

# Copy webui themes
COPY --from=builder /build /usr/share/amule/webserver

# copy local files
COPY root/ /

# Check binaries are OK
RUN chmod a+x /healthcheck.sh && \
    echo "**** checking binaries and themes ****" && \
    ls -l /usr/share/amule/webserver && \
    ldd /usr/bin/alcc && \
    ldd /usr/bin/amulecmd && \
    ldd /usr/bin/amuled && \
    ldd /usr/bin/amuleweb && \
    ldd /usr/bin/ed2k

# Define healthcheck script
HEALTHCHECK --start-period=2s --interval=30s --timeout=3s CMD /healthcheck.sh

# ports and volumes
EXPOSE "${PORT}" "${ECPORT}"
VOLUME "${DATADIR}"
