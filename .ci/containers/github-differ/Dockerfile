from golang:1.19-bullseye as resource
SHELL ["/bin/bash", "-c"]

RUN apt-get update
RUN apt-get install -y curl jq
ADD generate_comment.sh /generate_comment.sh
ADD compare_breaking_changes.sh /compare_breaking_changes.sh
ADD test_tools.sh /test_tools.sh
ADD utils.sh /utils.sh
<<<<<<< HEAD
RUN chmod +x /utils.sh
=======
>>>>>>> 31547b00e1e93da8fd198d133a17df4264d4eb4d
ENTRYPOINT ["/generate_comment.sh"]
