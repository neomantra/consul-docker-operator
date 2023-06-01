# consul-docker-operator Dockerfile
# Copyright (c) 2023 Neomantra BV

ARG BUILD_BASE="golang"
ARG BUILD_TAG="1.20-bullseye"

ARG RUNTIME_BASE="debian"
ARG RUNTIME_TAG="bullseye-slim"

ARG BUILD_ARCH="amd64"

##################################################################################################
# Builder
##################################################################################################

FROM ${BUILD_BASE}:${BUILD_TAG} AS build

ARG BUILD_BASE="golang"
ARG BUILD_TAG="1.20-bullseye"

ARG BUILD_ARCH="amd64"

RUN DEBIAN_FRONTEND=noninteractive apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    git

ARG TASKFILE_VERSION="v3.23.0"
RUN curl -fSL "https://github.com/go-task/task/releases/download/${TASKFILE_VERSION}/task_linux_${BUILD_ARCH}.deb" -o /tmp/task_linux.deb \
    && dpkg -i /tmp/task_linux.deb \
    && rm /tmp/task_linux.deb

RUN env
ADD . /src
WORKDIR /src

# Binaries
RUN mkdir bin && task


##################################################################################################
# Runtime environment
##################################################################################################

FROM ${RUNTIME_BASE}:${RUNTIME_TAG} AS runtime

ARG BUILD_BASE="golang"
ARG BUILD_TAG="1.20-bullseye"

ARG RUNTIME_BASE="debian"
ARG RUNTIME_TAG="bullseye-slim"

ARG BUILD_ARCH="amd64"

# Install dependencies and ops tools
RUN DEBIAN_FRONTEND=noninteractive apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
        ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy binaries
COPY --from=build /src/bin/* /usr/local/bin/

LABEL BUILD_BASE="${BUILD_BASE}"
LABEL BUILD_TAG="${BUILD_TAG}"
LABEL BUILD_ARCH="${BUILD_ARCH}"
LABEL RUNTIME_BASE="${RUNTIME_BASE}"
LABEL RUNTIME_TAG="${RUNTIME_TAG}"
