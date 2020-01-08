FROM alpine/git

RUN apk add --no-cache curl
RUN apk add --no-cache curl-dev
RUN apk add --no-cache bash
RUN apk add --no-cache jq
ADD generate_comment.sh /generate_comment.sh
ENTRYPOINT ["/generate_comment.sh"]
