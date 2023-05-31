from golang:1.19-bullseye as resource
SHELL ["/bin/bash", "-c"]

RUN apt-get update
RUN apt-get install -y git jq

ADD test_terraform.sh /test_terraform.sh
ENTRYPOINT ["/test_terraform.sh"]
