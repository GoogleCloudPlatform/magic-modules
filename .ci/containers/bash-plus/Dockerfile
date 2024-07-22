FROM alpine/git

RUN apk update && apk upgrade && \
    apk add --no-cache bash jq curl && \
    rm -rf /var/cache/apk/*
