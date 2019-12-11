from golang:1.11-stretch as resource

SHELL ["/bin/bash", "-c"]

RUN go get golang.org/x/tools/cmd/goimports

# Set up Github SSH cloning.
RUN ssh-keyscan github.com >> /known_hosts
RUN echo "UserKnownHostsFile /known_hosts" >> /etc/ssh/ssh_config

ENV GOFLAGS "-mod=vendor"

# Install Ruby from source.
RUN apt-get update
RUN apt-get install -y bzip2 libssl-dev libreadline-dev zlib1g-dev
RUN git clone https://github.com/rbenv/rbenv.git /rbenv
ENV PATH /rbenv/bin:/root/.rbenv/shims:$PATH

ENV RUBY_VERSION 2.6.0
ENV RUBYGEMS_VERSION 3.0.2
ENV BUNDLER_VERSION 1.17.0

RUN /rbenv/bin/rbenv init || true
RUN eval "$(rbenv init -)"
RUN mkdir -p "$(rbenv root)"/plugins
RUN git clone https://github.com/rbenv/ruby-build.git "$(rbenv root)"/plugins/ruby-build

RUN rbenv install $RUBY_VERSION
RUN rbenv global 2.6.0
RUN rbenv rehash

RUN gem update --system "$RUBYGEMS_VERSION"
RUN gem install bundler --version "$BUNDLER_VERSION" --force

ADD Gemfile Gemfile
ADD Gemfile.lock Gemfile.lock
RUN bundle install
RUN rbenv rehash
