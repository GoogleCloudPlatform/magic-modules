# Stage 1: Download go module cache for builds
FROM golang:1.19-bullseye AS builder

WORKDIR /app1
ADD "https://raw.githubusercontent.com/GoogleCloudPlatform/magic-modules/main/.ci/magician/go.mod" go.mod
ADD "https://raw.githubusercontent.com/GoogleCloudPlatform/magic-modules/main/.ci/magician/go.sum" go.sum
RUN go mod download

# Stage 2: Creating the final imag
FROM golang:1.19-bullseye
SHELL ["/bin/bash", "-c"]

# Copy Go dependencies from builder
COPY --from=builder /go/pkg/mod /go/pkg/mod

# Set up Github SSH cloning.
RUN ssh-keyscan github.com >> /known_hosts
RUN echo "UserKnownHostsFile /known_hosts" >> /etc/ssh/ssh_config

RUN apt-get update && \
    apt-get install -y git jq unzip parallel curl && \
    echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] http://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && \
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg  add - && \
    apt-get update -y && \
    apt-get install google-cloud-sdk -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN wget https://releases.hashicorp.com/terraform/1.2.5/terraform_1.2.5_linux_amd64.zip \
    && unzip terraform_1.2.5_linux_amd64.zip \
    && rm terraform_1.2.5_linux_amd64.zip \
    && mv ./terraform /bin/terraform

### the following version of terraform are required for
### terraform-validator-tester-integration
RUN wget https://releases.hashicorp.com/terraform/0.13.7/terraform_0.13.7_linux_amd64.zip \
    && unzip terraform_0.13.7_linux_amd64.zip \
    && rm terraform_0.13.7_linux_amd64.zip \
    && mkdir /terraform \
    && mv terraform /terraform/0.13.7

RUN wget https://releases.hashicorp.com/terraform/0.12.31/terraform_0.12.31_linux_amd64.zip \
    && unzip terraform_0.12.31_linux_amd64.zip \
    && rm terraform_0.12.31_linux_amd64.zip \
    && mv terraform /terraform/0.12.31