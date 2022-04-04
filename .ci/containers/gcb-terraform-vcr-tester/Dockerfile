from golang:1.16-stretch as resource
SHELL ["/bin/bash", "-c"]
# Set up Github SSH cloning.
RUN ssh-keyscan github.com >> /known_hosts
RUN echo "UserKnownHostsFile /known_hosts" >> /etc/ssh/ssh_config

RUN apt-get update
RUN apt-get install -y git jq

RUN curl -sSL https://sdk.cloud.google.com | bash
ENV PATH $PATH:/root/google-cloud-sdk/bin

ADD test_terraform_vcr.sh /test_terraform_vcr.sh
ENTRYPOINT ["/test_terraform_vcr.sh"]