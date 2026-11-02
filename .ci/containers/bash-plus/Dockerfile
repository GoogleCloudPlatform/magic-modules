FROM alpine/git:2.54.0@sha256:ce84fd6a057b4dc5cf5009da2dd832f1c6d096c1984bf50496cf3c256402264d

RUN apk update && apk upgrade && \
    apk add --no-cache bash jq curl && \
    rm -rf /var/cache/apk/*
