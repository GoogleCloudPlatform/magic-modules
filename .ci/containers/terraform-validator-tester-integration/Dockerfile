from golang:1.19-bullseye as resource
SHELL ["/bin/bash", "-c"]
# Set up Github SSH cloning.
RUN ssh-keyscan github.com >> /known_hosts
RUN echo "UserKnownHostsFile /known_hosts" >> /etc/ssh/ssh_config

RUN apt-get update
RUN apt-get install -y git jq unzip

ADD test_terraform_validator_integration.sh /test_terraform_validator_integration.sh
ENTRYPOINT ["/test_terraform_validator_integration.sh"]

RUN wget https://releases.hashicorp.com/terraform/0.13.7/terraform_0.13.7_linux_amd64.zip
RUN unzip terraform_0.13.7_linux_amd64.zip && rm terraform_0.13.7_linux_amd64.zip
RUN mkdir /terraform
RUN mv terraform /terraform/0.13.7

RUN wget https://releases.hashicorp.com/terraform/0.12.31/terraform_0.12.31_linux_amd64.zip
RUN unzip terraform_0.12.31_linux_amd64.zip && rm terraform_0.12.31_linux_amd64.zip
RUN mv terraform /terraform/0.12.31