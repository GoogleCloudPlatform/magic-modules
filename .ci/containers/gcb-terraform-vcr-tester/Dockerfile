from golang:1.16-stretch as resource
SHELL ["/bin/bash", "-c"]
# Set up Github SSH cloning.
RUN ssh-keyscan github.com >> /known_hosts
RUN echo "UserKnownHostsFile /known_hosts" >> /etc/ssh/ssh_config

RUN apt-get update
RUN apt-get install -y git jq unzip

RUN curl -sSL https://sdk.cloud.google.com | bash
ENV PATH $PATH:/root/google-cloud-sdk/bin

ADD test_terraform_vcr.sh /test_terraform_vcr.sh

RUN wget https://releases.hashicorp.com/terraform/0.14.8/terraform_0.14.8_linux_amd64.zip
RUN unzip terraform_0.14.8_linux_amd64.zip && rm terraform_0.14.8_linux_amd64.zip
RUN mv ./terraform /bin/terraform

ENTRYPOINT ["/test_terraform_vcr.sh"]