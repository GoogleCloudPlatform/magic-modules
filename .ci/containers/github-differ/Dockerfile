from golang:1.19-bullseye as resource
SHELL ["/bin/bash", "-c"]

RUN apt-get update
RUN apt-get install -y curl jq yamllint
ADD generate_comment.sh /generate_comment.sh
ADD compare_breaking_changes.sh /compare_breaking_changes.sh
ADD test_tools.sh /test_tools.sh
ADD utils.sh /utils.sh
ENTRYPOINT ["/generate_comment.sh"]
