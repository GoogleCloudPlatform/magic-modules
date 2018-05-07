FROM ruby:2.5

ENV PUPPET_AGENT_VERSION="5.5.0"

# Install puppet agent
RUN apt-get update && \
	wget https://apt.puppetlabs.com/puppet5-release-stretch.deb && \
	dpkg -i puppet5-release-stretch.deb && \
	rm puppet5-release-stretch.deb && \
	apt-get update && \
    apt-get install --no-install-recommends -y puppet-agent="$PUPPET_AGENT_VERSION"-1stretch && \
    apt-get clean

# Install gem required by gauth
RUN gem install googleauth google-api-client

# Setup symlink for puppet modules
COPY docker-entrypoint.sh /usr/local/bin

ENV PATH=/opt/puppetlabs/server/bin:/opt/puppetlabs/puppet/bin:/opt/puppetlabs/bin:$PATH
ENV FACTER_cred_path=/etc/creds.json

WORKDIR ~

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["/bin/bash"]