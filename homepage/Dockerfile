# syntax=docker/dockerfile:1
FROM node:22-alpine3.22 AS builder

ARG BUILD_DATE="unknown"
ARG VERSION="latest"
ARG COMMIT_BRANCH="main"

RUN echo "**** install runtime packages ****" && \
  apk add --no-cache git unzip wget xz bash ca-certificates curl

COPY patches/ /patches/

RUN echo "**** Download homepage and apply patches ****" && \
  git clone --depth 1 --single-branch --branch $COMMIT_BRANCH https://github.com/gethomepage/homepage.git /app && \
  cd /app && for patch in /patches/*.patch; do echo "Applying patch: $patch"; git apply "$patch"; done

WORKDIR /app

ENV NEXT_TELEMETRY_DISABLED=1
ENV NEXT_PUBLIC_BUILDTIME=$BUILD_DATE
ENV NEXT_PUBLIC_VERSION=$VERSION
ENV NEXT_PUBLIC_REVISION=$COMMIT_BRANCH

RUN echo "**** Install homepage dependencies and build ****" && \
  corepack enable && corepack prepare pnpm@latest --activate && \
  pnpm install --frozen-lockfile --prefer-offline && \
  pnpm run build

### Final stage
FROM ghcr.io/linuxserver/baseimage-alpine:3.22

# set version label
ARG BUILD_DATE="unknown"
ARG VERSION="latest"
ARG COMMIT_BRANCH="main"

LABEL org.opencontainers.image.created="${BUILD_DATE}"
LABEL org.opencontainers.image.authors="linuxserver.io authors and jriguera"
LABEL org.opencontainers.image.url="https://github.com/jriguera/container-images/tree/main/homepage"
LABEL org.opencontainers.image.documentation="https://github.com/jriguera/container-images/tree/main/homepage"
LABEL org.opencontainers.image.source="https://github.com/jriguera/container-images"
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.revision="${COMMIT_BRANCH}"
LABEL org.opencontainers.image.vendor="jriguera"
LABEL org.opencontainers.image.licenses="GPL-3.0-only"
LABEL org.opencontainers.image.ref.name="${COMMIT_BRANCH}"
LABEL org.opencontainers.image.title="Homepage"
LABEL org.opencontainers.image.description="[Homepage](https://gethomepage.dev/) is a self-hosted dashboard for your home server."
LABEL build_version="Homepage version:- ${VERSION} Build-date:- ${BUILD_DATE}"
LABEL maintainer="jriguera"

# environment variables
ENV CONFIGDIR=/config
ENV HOMEPAGE_CONFIG_DIR=/config
ENV NODE_ENV=production
ENV HOSTNAME=0.0.0.0
ENV PORT=3000

# copy local files
COPY root/ /

# Install packages
RUN echo "**** install runtime packages ****" && \
  apk add --no-cache \
    tzdata \
    pwgen \
    curl \
    coreutils \
    bash \
    socat \
    net-tools \
    iputils-ping \
    shadow \
    npm \
    nodejs \
  && \
  printf "Homepage version: ${VERSION}\nBuild-date: ${BUILD_DATE}" > /build_version && \
  echo "**** cleanup ****" && \
  chmod a+x /healthcheck.sh && \
  rm -rf /tmp/* $HOME/.cache

# Copy only necessary files from the build stage
COPY --from=builder --chown=abc:abc /app/public/ /public
COPY --from=builder --chown=abc:abc /app/.next/standalone /app
COPY --from=builder --chown=abc:abc /app/.next/static/ /app/.next/static

RUN echo "**** create symlinks ****" && cd /app && \
  ln -s /config config && \
  ln -s /public public

# Define healthcheck script
HEALTHCHECK --start-period=2s --interval=30s --timeout=3s CMD /healthcheck.sh

# ports and volumes
EXPOSE "${PORT}"
VOLUME "${CONFIGDIR}"
