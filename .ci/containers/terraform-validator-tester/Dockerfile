from golang:1.19-stretch as resource
SHELL ["/bin/bash", "-c"]
# Set up Github SSH cloning.
RUN ssh-keyscan github.com >> /known_hosts
RUN echo "UserKnownHostsFile /known_hosts" >> /etc/ssh/ssh_config

RUN apt-get update
RUN apt-get install -y git jq

ADD test_terraform_validator.sh /test_terraform_validator.sh
ENTRYPOINT ["/test_terraform_validator.sh"]