from golang:1.13-stretch as resource

SHELL ["/bin/bash", "-c"]

RUN go get golang.org/x/tools/cmd/goimports

# Set up Github SSH cloning.
RUN ssh-keyscan github.com >> /known_hosts
RUN echo "UserKnownHostsFile /known_hosts" >> /etc/ssh/ssh_config

ENV GOFLAGS "-mod=vendor"
ENV GO111MODULE "off"

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
RUN rbenv global $RUBY_VERSION
RUN rbenv rehash

RUN gem update --system "$RUBYGEMS_VERSION"
RUN gem install bundler --version "$BUNDLER_VERSION" --force

RUN apt-get update
RUN apt-get install -y git build-essential libbz2-dev libssl-dev libreadline-dev \
                        libffi-dev libsqlite3-dev tk-dev libpng-dev libfreetype6-dev \
                        make build-essential libssl-dev zlib1g-dev libbz2-dev \
                        libreadline-dev libsqlite3-dev curl llvm libncurses5-dev libncursesw5-dev \
                        xz-utils tk-dev libffi-dev liblzma-dev python-openssl

RUN wget https://raw.githubusercontent.com/GoogleCloudPlatform/magic-modules/master/Gemfile
RUN wget https://raw.githubusercontent.com/GoogleCloudPlatform/magic-modules/master/Gemfile.lock
RUN bundle install
RUN rbenv rehash
RUN rm Gemfile Gemfile.lock

RUN curl -L https://github.com/pyenv/pyenv-installer/raw/master/bin/pyenv-installer | bash
ENV PATH="/root/.pyenv/bin:${PATH}"
RUN eval "$(pyenv init -)"
RUN eval "$(pyenv virtualenv-init -)"
RUN pyenv install 3.6.8
RUN pyenv install 2.7.13
RUN pyenv rehash
ENV PATH="/root/.pyenv/shims:${PATH}"
RUN pyenv global 2.7.13 3.6.8
RUN pip install beautifulsoup4 mistune
RUN pip3 install black
ENV LC_ALL=C.UTF-8
ENV LANG=C.UTF-8

RUN git config --global user.name "Modular Magician"
RUN git config --global user.email "magic-modules@google.com"

ADD generate_downstream.sh /generate_downstream.sh

ENTRYPOINT ["/generate_downstream.sh"]
