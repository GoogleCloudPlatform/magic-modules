FROM alpine/git

RUN apk add --no-cache curl
RUN apk add --no-cache bash
ADD generate_comment.sh /generate_comment.sh
ENTRYPOINT ["/generate_comment.sh"]
