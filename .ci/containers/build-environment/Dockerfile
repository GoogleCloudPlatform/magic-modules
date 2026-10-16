FROM golang:1.26.9-trixie@sha256:f06b7122f888680048afb247c4e5f49c0069eb6c4c1ccae0f3f4834527a8d22c AS builder

# Set working directory
WORKDIR /app

# Download go.mod and go.sum from your repository
ADD "https://raw.githubusercontent.com/GoogleCloudPlatform/magic-modules/main/tpgtools/go.mod" go.mod
ADD "https://raw.githubusercontent.com/GoogleCloudPlatform/magic-modules/main/tpgtools/go.sum" go.sum

# Install the go dependencies
RUN go mod download

ENV GOPATH /go
ENV PATH /usr/local/go/bin:$PATH
ENV PATH $GOPATH/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 1777 "$GOPATH"
WORKDIR $GOPATH

# terraform binary used by tfv/tgc
COPY --from=hashicorp/terraform:1.14.3@sha256:20056ea0d2706f2a944e202befa8a558dbcc5798da0b64f884b2647aab82a4ef /bin/terraform /bin/terraform

SHELL ["/bin/bash", "-c"]

ENV GO111MODULE "on"
ENV LC_ALL=C.UTF-8
ENV LANG=C.UTF-8

RUN apt-get update && apt-get install -y --no-install-recommends git openssh-client apt-transport-https ca-certificates curl netbase wget gcc make jq libjq1 unzip zip

RUN git config --global user.name "Modular Magician"
RUN git config --global user.email "magic-modules@google.com"

RUN go install golang.org/x/tools/cmd/goimports@d088b475e3360caabc032aaee1dc66351d4e729a
RUN go install github.com/github/hub@dcdf56167c38eacef96ad5ef88575e7f74f9a6f2 # v2.11.2+incompatible
