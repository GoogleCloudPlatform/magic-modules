#!/usr/bin/env bash
# run-in-container runs a make command within a
# container. The only required dependency is
# podman or docker: containers required for runtime
# will be pulled or built during the exectuion of the script.
set -e

BUILDER_IMAGE="gcr.io/graphite-docker-images/downstream-builder:latest"
DEV_IMAGE="gcr.io/graphite-docker-images/magic-modules-dev:latest"
GOPATH="/Users/enkuwfeleke/go"

main() {
    if which podman > /dev/null; then
        CONTAINER_EXECUTABLE="podman"
    elif which docker > /dev/null; then
        CONTAINER_EXECUTABLE="docker"
    else
        echo "Unable to find podman or docker executable. Please install podman or docker and try again." && exit 1;
    fi

    echo "Found ${CONTAINER_EXECUTABLE}, using it to run commands."

    if ! ${CONTAINER_EXECUTABLE} images | grep "gcr.io/graphite-docker-images/magic-modules-dev" > /dev/null; then
        echo "unable to find magic-modules-dev container"
        build_dev_container
    fi
    echo $GOPATH
    ${CONTAINER_EXECUTABLE} run --rm \
    -v "$PWD:/magic-modules" -v "$GOPATH:$GOPATH" \
    -it "${DEV_IMAGE}" "$@"
}


build_dev_container() {
    echo "building containers...."
    echo "NOTE: this may take a while."
    sleep 1
    if ! ${CONTAINER_EXECUTABLE} images | grep "${BUILDER_IMAGE}" > /dev/null; then
        if ! ${CONTAINER_EXECUTABLE} pull "${BUILDER_IMAGE}"; then
            ${CONTAINER_EXECUTABLE} build \
                .ci/containers/downstream-builder/ \
                -t "${BUILDER_IMAGE}"
        fi
    fi
    ${CONTAINER_EXECUTABLE} build \
        .ci/containers/magic-modules-dev/ \
        -t "${DEV_IMAGE}"
}

main "$@"