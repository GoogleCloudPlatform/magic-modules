from golang:1.19-stretch as resource
SHELL ["/bin/bash", "-c"]
# Set up Github SSH cloning.
RUN ssh-keyscan github.com >> /known_hosts
RUN echo "UserKnownHostsFile /known_hosts" >> /etc/ssh/ssh_config

RUN apt-get update
RUN apt-get install -y git jq

ADD teamcityparams.xml /teamcityparams.xml
ADD run_vcr_tests.sh /run_vcr_tests.sh
ADD check_membership.sh /check_membership.sh
ADD teamcityparamsrecording.xml /teamcityparamsrecording.xml
ADD teamcitycancelparams.xml /teamcitycancelparams.xml
ENTRYPOINT ["/check_membership.sh"]
