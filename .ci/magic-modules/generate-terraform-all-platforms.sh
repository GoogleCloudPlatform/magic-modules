#!/usr/bin/env bash


set -x

function init {
  START_DIR=${PWD}
  # Setup GOPATH
  export GOPATH=${PWD}/go
  # Setup GOBIN
  export GOBIN=${PWD}/dist
  # Create GOBIN folder
  mkdir -p "$GOBIN"
  # Create GOPATH structure
  mkdir -p "${GOPATH}/src/github.com/terraform-providers"
  # Copy the repo
  cp -rf "$1" "${GOPATH}/src/github.com/terraform-providers/terraform-provider-google"
  # Paths and vars
  PROVIDER_NAME="google"
  PROVIDERPATH="$GOPATH/src/github.com/terraform-providers"
  SRC_DIR="$PROVIDERPATH/terraform-provider-$PROVIDER_NAME"
  TARGET_DIR="$START_DIR/dist"
  XC_ARCH=${XC_ARCH:-"386 amd64 arm"}
  XC_OS=${XC_OS:=linux darwin windows freebsd openbsd solaris}
  XC_EXCLUDE_OSARCH="!darwin/arm !darwin/386"
  export CGO_ENABLED=0
  mkdir -p "$TARGET_DIR"
}

function installGox {
  if ! which gox > /dev/null; then
    go get -u github.com/mitchellh/gox
  fi
}

function compile {
  pushd "$SRC_DIR"
  printf "\n"
  make fmtcheck

  # Set LD Flags
  LD_FLAGS="-s -w"

  # Clean any old directories (should never be here)
  rm -f bin/*
  rm -fr pkg/*
  # Build with gox
  "$GOBIN/gox" \
    -os="${XC_OS}" \
    -arch="${XC_ARCH}" \
    -osarch="${XC_EXCLUDE_OSARCH}" \
    -ldflags "${LD_FLAGS}" \
    -output "$TARGET_DIR/terraform-provider-${PROVIDER_NAME}.{{.OS}}_{{.Arch}}" \
    .

  popd
}

function main {
  init "$1"
  installGox
  compile
}

main "$@"
