from golang:1.19-stretch as resource
SHELL ["/bin/bash", "-c"]

RUN apt-get update
RUN apt-get install -y curl jq
ADD generate_comment.sh /generate_comment.sh
ADD compare_breaking_changes.sh /compare_breaking_changes.sh
ADD test_tools.sh /test_tools.sh
ENTRYPOINT ["/generate_comment.sh"]
