# Stage 1: Download go module cache for builds
FROM golang:1.23-bullseye AS builder
ENV GOCACHE=/go/cache

RUN apt-get update && apt-get install -y unzip

# Download tpgtools dependencies (from build-env)
WORKDIR /app
ADD "https://raw.githubusercontent.com/GoogleCloudPlatform/magic-modules/main/tpgtools/go.mod" go.mod
ADD "https://raw.githubusercontent.com/GoogleCloudPlatform/magic-modules/main/tpgtools/go.sum" go.sum
RUN go mod download

WORKDIR /app1
# Add the source code and build
ADD "https://github.com/GoogleCloudPlatform/magic-modules/archive/refs/heads/main.zip" source.zip
RUN unzip source.zip && rm source.zip
WORKDIR /app1/magic-modules-main/.ci/magician
# Build the binary (we won't use it in the final image, but it's cached)
RUN go build -o /dev/null .

# Stage 2: Creating the final image
FROM golang:1.23-bullseye
SHELL ["/bin/bash", "-c"]
ENV GOCACHE=/go/cache

# Copy Go dependencies and Go build cache
COPY --from=builder /go/pkg/mod /go/pkg/mod
COPY --from=builder /go/cache /go/cache

# Add build-env environment variables
ENV GOPATH /go
ENV PATH /usr/local/go/bin:$PATH
ENV PATH $GOPATH/bin:$PATH
ENV GO111MODULE "on"
ENV LC_ALL=C.UTF-8
ENV LANG=C.UTF-8

# Create GOPATH structure (from build-env)
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 1777 "$GOPATH"

RUN apt-get update && \
    apt-get install -y git jq unzip parallel curl && \
    # Add build-env packages
    apt-get install -y openssh-client apt-transport-https ca-certificates netbase wget gcc make libjq1 && \
    echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] http://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && \
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg  add - && \
    apt-get update -y && \
    apt-get install google-cloud-sdk -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Add git configuration (from build-env)
RUN git config --global user.name "Modular Magician"
RUN git config --global user.email "magic-modules@google.com"

RUN wget https://releases.hashicorp.com/terraform/1.11.0/terraform_1.11.0_linux_amd64.zip \
    && unzip terraform_1.11.0_linux_amd64.zip \
    && rm terraform_1.11.0_linux_amd64.zip \
    && mv ./terraform /bin/terraform

# Install Go tools (from build-env)
RUN go install golang.org/x/tools/cmd/goimports@d088b475e3360caabc032aaee1dc66351d4e729a
RUN go install github.com/github/hub@v2.11.2+incompatible