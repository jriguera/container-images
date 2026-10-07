# syntax=docker/dockerfile:1

ARG AMULE_VERSION=3.1.0
ARG BUILD_DATE="unknown"
ARG VERSION="latest"
ARG COMMIT_SHA="main"
ARG COMMIT_BRANCH="main"

FROM alpine:edge AS builder

ARG AMULE_VERSION
ARG COMMIT_BRANCH

# Install packages
RUN \
  echo "**** install runtime packages ****" && \
  apk add --no-cache git unzip wget xz bash ca-certificates curl

RUN \
  echo "**** install build packages ****" && \
  apk add --no-cache \
    build-base \
    cmake \
    make \
    git \
    binutils-dev \
    ca-certificates \
    pkgconf \
    boost-dev \
    libmaxminddb-dev \
    glib-dev \
    readline-dev \
#   wxwidgets-dev \
    zlib-dev \
    libpng-dev \
    libupnp-dev \
    curl-dev \
    crypto++-dev

WORKDIR /build

RUN \
  echo "**** Download webui themes (shallow clones) ****" && \
  git clone --depth 1 --single-branch https://github.com/pedro77/amuleweb-bootstrap-template.git bootstrap.git && \
  mv bootstrap.git/dist /build/bootstrap && rm -rf bootstrap.git && \
  git clone --depth 1 --single-branch https://github.com/MatteoRagni/AmuleWebUI-Reloaded reloaded && \
  rm -rf reloaded/.git* reloaded/doc-images && \
  git clone --depth 1 --single-branch https://github.com/esaracho/amuleweb-adaptable adaptable && \
  rm -rf adaptable/.git*

# wxWidgets is built from source because of a bug in wx: amuled hits a use-after-free in
# the wx event loop during the startup HTTP downloads
# https://github.com/amule-org/amule/issues/1136
# The fix lands in wx 3.2.12, so WX_REF pins the 3.2-branch commit that carries it.
# TODO: revert this and install the distro wxwidgets-dev/wxwidgets packages again once
# wx 3.2.12 is released/packaged
# Only wxBase and wxNet are built, which also drops glib: aMule links it only when wx
# reports __WXGTK__. --with-libcurl is mandatory, aMule's HTTP downloads use wxWebRequest
ARG WX_REF=3382306c8c35113ea93fe9276aa6f3a0587677c0
# git clone --branch only accepts branch/tag names, so fetch the commit SHA directly
RUN git init -q /build/wx-src && cd /build/wx-src && \
    git fetch -q --depth 1 https://github.com/wxWidgets/wxWidgets.git "${WX_REF}" && git checkout -q FETCH_HEAD && \
    ./configure \
        --prefix=/usr \
        --disable-gui \
        --disable-debug \
        --disable-precomp-headers \
        --without-subdirs \
        --without-expat \
        --with-regex=sys \
        --with-zlib=sys \
        --with-libcurl \
        --disable-mimetype \
        --disable-secretstore \
        --disable-debugreport \
        --disable-fs_inet \
        --disable-ftp \
        --disable-sound && \
    make -j"$(nproc)" && \
    make install && \
    strip /usr/lib/libwx_baseu*.so.*.*.* && \
    cd /build && rm -rf wx-src

RUN \
  echo "**** Download amule ****" && \
  git clone --depth 1 --single-branch --branch "$AMULE_VERSION" https://github.com/amule-org/amule.git /build/amule

RUN echo "**** Build amule ****" && \
  cmake -B amule-build amule \
        -DCMAKE_INSTALL_PREFIX=/usr \
        -DCMAKE_BUILD_TYPE=Release \
        -DBUILD_TESTING=NO \
        -DBUILD_MONOLITHIC=NO \
        -DBUILD_DAEMON=YES \
        -DBUILD_AMULECMD=YES \
        -DBUILD_WEBSERVER=YES \
        -DBUILD_AMULEAPI=YES \
        -DBUILD_ALCC=YES \
        -DBUILD_FILEVIEW=YES \
        -DENABLE_IP2COUNTRY=YES \
        -DENABLE_UPNP=YES \
        -DENABLE_NLS=NO && \
  cmake --build amule-build -j"$(nproc)" && \
  cmake --install amule-build

FROM ghcr.io/linuxserver/baseimage-alpine:3.24

# set version label
ARG BUILD_DATE
ARG VERSION
ARG COMMIT_SHA
ARG AMULE_VERSION

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
ENV PORT_WEB=4711
ENV ECPORT=4712
ENV AMULE_USER=abc
ENV ENABLE_AMULEWEB=false
ENV ENABLE_AMULEAPI=true

# Use the built-in C.UTF-8 locale so amuled handles non-ASCII paths correctly
ENV LANG=C.UTF-8
ENV LC_ALL=C.UTF-8

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
    crypto++ \
    ffmpeg \
    binutils \
    libmaxminddb \
  && \
  printf "Amule version: ${VERSION}\nBuild-date: ${BUILD_DATE}" > /build_version && \
  echo "**** cleanup ****" && \
  rm -rf /tmp/* $HOME/.cache

# Copy compiled binaries
COPY --from=builder /usr/bin/alcc /usr/bin/amuleapi /usr/bin/amulecmd \
     /usr/bin/amuled /usr/bin/amuleweb /usr/bin/ed2k /usr/bin/fileview /usr/bin/

# Copy amule shared data (includes default webserver themes)
COPY --from=builder /usr/share/amule /usr/share/amule

# Copy extra webui themes (bootstrap, reloaded, adaptable)
COPY --from=builder /build/bootstrap /usr/share/amule/webserver/bootstrap
COPY --from=builder /build/reloaded /usr/share/amule/webserver/reloaded
COPY --from=builder /build/adaptable /usr/share/amule/webserver/adaptable

# copy local files
COPY root/ /

# Check binaries are OK
RUN chmod a+x /healthcheck.sh && \
    echo "**** checking binaries and themes ****" && \
    ls -l /usr/share/amule/webserver && \
    ldd /usr/bin/alcc && \
    ldd /usr/bin/amuleapi && \
    ldd /usr/bin/amulecmd && \
    ldd /usr/bin/amuled && \
    ldd /usr/bin/amuleweb && \
    ldd /usr/bin/ed2k && \
    ldd /usr/bin/fileview && \
    ldd "$(command -v ffprobe)"

# Define healthcheck script
HEALTHCHECK --start-period=2s --interval=30s --timeout=3s CMD /healthcheck.sh

# ports and volumes
EXPOSE "${PORT_WEB}" "${ECPORT}"
VOLUME "${DATADIR}"
