# Stage 1: Download go module cache for builds
FROM golang:1.25-trixie AS builder
ENV GOCACHE=/go/cache

RUN apt-get update && apt-get install -y unzip
WORKDIR /app1
# Add the source code and build
ADD "https://github.com/GoogleCloudPlatform/magic-modules/archive/refs/heads/main.zip" source.zip
RUN unzip source.zip && rm source.zip
WORKDIR /app1/magic-modules-main/.ci/magician
# Build the binary (we won't use it in the final image, but it's cached)
RUN go build -o /dev/null .

# Stage 2: Creating the final image
FROM golang:1.25-trixie
SHELL ["/bin/bash", "-c"]
ENV GOCACHE=/go/cache

# Copy Go dependencies and Go build cache
COPY --from=builder /go/pkg/mod /go/pkg/mod
COPY --from=builder /go/cache /go/cache

RUN apt-get update && \
    apt-get install -y git jq unzip zip parallel apt-transport-https ca-certificates gnupg curl && \
    echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && \
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | gpg --dearmor -o /usr/share/keyrings/cloud.google.gpg && \
    apt-get update -y && \
    apt-get install google-cloud-cli -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN wget https://releases.hashicorp.com/terraform/1.14.3/terraform_1.14.3_linux_amd64.zip \
    && unzip terraform_1.14.3_linux_amd64.zip \
    && rm terraform_1.14.3_linux_amd64.zip \
    && mv ./terraform /bin/terraform
