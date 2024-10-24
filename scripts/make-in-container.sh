#!/usr/bin/env bash
# run-in-container runs a make command within a
# container. The only required dependency is
# podman or docker: containers required for runtime
# will be pulled or built during the exectuion of the script.
set -e

BUILDER_IMAGE="gcr.io/graphite-docker-images/build-environment:testing"

if [ -t 1 ]; then
    red=$'\e[1;31m'
    grn=$'\e[1;32m'
    yel=$'\e[1;33m'
    blu=$'\e[1;34m'
    mag=$'\e[1;35m'
    cyn=$'\e[1;36m'
    end=$'\e[0m'
else
    red=''
    grn=''
    yel=''
    blu=''
    mag=''
    cyn=''
    end=''
fi

warn() {
    printf "%s\n" "${!1}${2}${end}"
}

main() {
    exitcode=0

    if [[ -z "${GOPATH}" ]]; then
        warn "red" "GOPATH is not set on the host machine. Add 'export GOPATH=\$HOME/go' to your terminal's startup script and restart your terminal."
        exitcode=1
    else
        echo "Host GOPATH=${GOPATH}"
    fi


    if which docker > /dev/null || which podman > /dev/null; then
        if which docker > /dev/null; then
            CONTAINER_EXECUTABLE="docker"
        else
            CONTAINER_EXECUTABLE="podman"
        fi
        echo "Found ${CONTAINER_EXECUTABLE}, using it to run commands."

        echo "Updating containers..."
        echo "NOTE: this may take a while."
        sleep 1
        if ! ${CONTAINER_EXECUTABLE} images | grep "${BUILDER_IMAGE}" > /dev/null; then
            if ! ${CONTAINER_EXECUTABLE} pull "${BUILDER_IMAGE}"; then
                ${CONTAINER_EXECUTABLE} build \
                    .ci/containers/build-environment/ \
                    -t "${BUILDER_IMAGE}"
            fi
        fi
    else
        warn "red" "Unable to find podman or docker executable. Please install podman or docker and try again.";
        exitcode=1
    fi

    if [ $exitcode -eq 0 ]; then
        ${CONTAINER_EXECUTABLE} run --rm \
            -v "$PWD:/magic-modules" -v "$GOPATH:$GOPATH" \
            --entrypoint "/usr/bin/make" \
            -w "/magic-modules" \
            -it "${BUILDER_IMAGE}" "$@"
    else
        warn "red" "Found errors; exiting"
    fi

    exit $exitcode
}

main "$@"