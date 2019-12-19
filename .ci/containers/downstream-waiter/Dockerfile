FROM alpine/git

RUN apk add --no-cache bash
ADD wait_for_commit.sh /wait_for_commit.sh
ENTRYPOINT ["/wait_for_commit.sh"]
