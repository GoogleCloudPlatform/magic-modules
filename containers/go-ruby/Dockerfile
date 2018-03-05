from golang:1.9

RUN ssh-keyscan github.com >> /known_hosts
RUN echo "UserKnownHostsFile /known_hosts" >> /etc/ssh/ssh_config
RUN apt-get update
RUN apt-get install -y bzip2 libssl-dev libreadline-dev zlib1g-dev

ENV RUBY_VERSION 2.5.0
ENV RUBYGEMS_VERSION 2.7.4
ENV BUNDLER_VERSION 1.16.1

RUN git clone https://github.com/rbenv/ruby-build.git
RUN PREFIX=/usr/local ./ruby-build/install.sh
RUN ruby-build "$RUBY_VERSION" /usr/

RUN gem update --system "$RUBYGEMS_VERSION"
RUN gem install bundler --version "$BUNDLER_VERSION" --force

# install things globally, for great justice
# and don't create ".bundle" in all our apps
ENV GEM_HOME /usr/local/bundle
ENV BUNDLE_PATH="$GEM_HOME" \
	BUNDLE_BIN="$GEM_HOME/bin" \
	BUNDLE_SILENCE_ROOT_WARNING=1 \
	BUNDLE_APP_CONFIG="$GEM_HOME"
ENV PATH $BUNDLE_BIN:$PATH
RUN mkdir -p "$GEM_HOME" "$BUNDLE_BIN" \
	&& chmod 777 "$GEM_HOME" "$BUNDLE_BIN"

