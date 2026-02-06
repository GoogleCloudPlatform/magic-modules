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

# STAGE 2: Final Image
FROM golang:1.25-bookworm

SHELL ["/bin/bash", "-c"]
ENV GOPATH=/go
ENV GOCACHE=/go/cache

# Copy pre-downloaded modules and build cache
COPY --from=builder /go/pkg/mod /go/pkg/mod
COPY --from=builder /go/cache /go/cache

# Install dependencies in a single layer to reduce image size
# Added --no-install-recommends to skip useless extras (docs, man pages)
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    git jq unzip zip parallel apt-transport-https ca-certificates gnupg curl && \
    echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && \
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | gpg --dearmor -o /usr/share/keyrings/cloud.google.gpg && \
    apt-get update -y && \
    apt-get install -y --no-install-recommends google-cloud-cli && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install Terraform
# Use a specific version variable for easier updates
ARG TF_VERSION=1.14.3
RUN curl -O https://releases.hashicorp.com/terraform/${TF_VERSION}/terraform_${TF_VERSION}_linux_amd64.zip \
    && unzip terraform_${TF_VERSION}_linux_amd64.zip \
    && rm terraform_${TF_VERSION}_linux_amd64.zip \
    && mv ./terraform /bin/terraform