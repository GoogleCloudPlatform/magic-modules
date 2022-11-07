from golang:1.18-stretch as resource
SHELL ["/bin/bash", "-c"]

RUN apt-get update
RUN apt-get install -y curl jq
ADD generate_comment.sh /generate_comment.sh
ADD compare_breaking_changes.sh /compare_breaking_changes.sh
ENTRYPOINT ["/generate_comment.sh"]
