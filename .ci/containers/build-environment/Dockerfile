# STAGE 1: Builder
FROM golang:1.25-bookworm AS builder

ENV GOPATH=/go
ENV GOCACHE=/go/cache
RUN mkdir -p /go/pkg/mod /go/cache

# Install minimal deps
RUN apt-get update && apt-get install -y unzip && rm -rf /var/lib/apt/lists/*

WORKDIR /app1

# NOTE: This ADD command busts the cache whenever the repo updates.
# There is no easy fix for this without changing how you supply the source code.
ADD "https://github.com/GoogleCloudPlatform/magic-modules/archive/refs/heads/main.zip" source.zip
RUN unzip -q source.zip && rm source.zip

# Silence the build output to keep logs clean, unless debugging
WORKDIR /app1/magic-modules-main/.ci/magician
RUN go build -o /dev/null .

# Consolidated dependencies download
# We group these to reduce layers, though it's mostly aesthetic in a multi-stage build
RUN cd /app1/magic-modules-main/mmv1 && go mod download -x && \
    cd /app1/magic-modules-main/tpgtools && go mod download -x && \
    cd /app1/magic-modules-main/mmv1/third_party/terraform && go mod download -x && \
    cd /app1/magic-modules-main/mmv1/third_party/tgc_next && go mod download -x

# Set up the Diff Processor directories
RUN mkdir -p /app1/magic-modules-main/tools/diff-processor/new /app1/magic-modules-main/tools/diff-processor/old && \
    cp /app1/magic-modules-main/mmv1/third_party/terraform/go.mod /app1/magic-modules-main/tools/diff-processor/new/go.mod && \
    cp /app1/magic-modules-main/mmv1/third_party/terraform/go.sum /app1/magic-modules-main/tools/diff-processor/new/go.sum && \
    cp /app1/magic-modules-main/mmv1/third_party/terraform/go.mod /app1/magic-modules-main/tools/diff-processor/old/go.mod && \
    cp /app1/magic-modules-main/mmv1/third_party/terraform/go.sum /app1/magic-modules-main/tools/diff-processor/old/go.sum && \
    cd /app1/magic-modules-main/tools/diff-processor/ && go mod download -x 

# Stage 2: Creating the final image
FROM golang:1.25-bookworm
SHELL ["/bin/bash", "-c"]

# Copy Go dependencies and Go build cache
COPY --from=builder /go/pkg/mod /go/pkg/mod
COPY --from=builder /go/cache /go/cache

ENV GOPATH /go
ENV GOCACHE=/go/cache
ENV PATH /usr/local/go/bin:$PATH
ENV PATH $GOPATH/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 1777 "$GOPATH"
WORKDIR $GOPATH

# terraform binary used by tfv/tgc
COPY --from=hashicorp/terraform:1.14.3 /bin/terraform /bin/terraform

SHELL ["/bin/bash", "-c"]

ENV GO111MODULE "on"
ENV LC_ALL=C.UTF-8
ENV LANG=C.UTF-8

RUN apt-get update && apt-get install -y --no-install-recommends git openssh-client apt-transport-https ca-certificates curl netbase wget gcc make jq libjq1 unzip zip

RUN git config --global user.name "Modular Magician"
RUN git config --global user.email "magic-modules@google.com"

RUN go install github.com/github/hub@v2.11.2+incompatible
