#!/usr/bin/env bash

set -e

if which podman > /dev/null; then
  CONTAINER_EXECUTABLE="podman"
elif which docker > /dev/null; then
  CONTAINER_EXECUTABLE="docker"
  ${CONTAINER_EXECUTABLE} buildx inspect amd64-builder || ${CONTAINER_EXECUTABLE} buildx create --name amd64-mybuilder --use --bootstrap
else
  echo "Unable to find podman or docker executable. Please install podman or docker and try again." && exit 1;
fi

${CONTAINER_EXECUTABLE} buildx build \
  --platform linux/amd64 \
  --output "type=tar,dest=./hello-cloud-run-grpc-container-image.tar" \
  --build-arg BUILDKIT_INLINE_CACHE=1 .
