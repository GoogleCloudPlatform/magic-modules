#!/usr/bin/env bash

set -e

# Setup GOPATH
export GOPATH=${PWD}/go
# Setup GOBIN
export GOBIN=${PWD}/dist

set -x

# Create GOBIN folder
mkdir -p "$GOBIN"

# Create GOPATH structure
mkdir -p "${GOPATH}/src/github.com/terraform-providers"
yes | cp -rf "$1" "${GOPATH}/src/github.com/terraform-providers/"

cd "${GOPATH}/src/github.com/terraform-providers/terraform-provider-google"

# Platforms build
platforms=("linux/386" "linux/amd64" "linux/arm" "linux/arm64" "linux/ppc64" "linux/ppc64le" "linux/mips" "linux/mipsle" "linux/mips64" "linux/mips64le")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    output_name='terraform-provider-google.'$GOOS'-'$GOARCH

    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -o $GOBIN/$output_name
done