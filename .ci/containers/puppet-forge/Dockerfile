from alpine:latest as resource

# Install Ruby from source.
RUN apk update
RUN apk add --no-cache --virtual .ruby-builddeps \
		autoconf \
		bash \
		bison \
		bzip2 \
		bzip2-dev \
		ca-certificates \
		coreutils \
		dpkg-dev dpkg \
		gcc \
		g++ \
		gdbm-dev \
		git \
		glib-dev \
		libc-dev \
		libffi-dev \
		libressl \
		libressl-dev \
		libxml2-dev \
		libxslt-dev \
		linux-headers \
		make \
		ncurses-dev \
		openssh-client \
		procps \
		readline-dev \
		ruby \
		tar \
		xz \
		yaml-dev \
		zlib-dev
ENV RUBY_VERSION 2.5.0
ENV RUBYGEMS_VERSION 2.7.4
ENV BUNDLER_VERSION 1.16.1

# Set up Github SSH cloning.
RUN ssh-keyscan github.com >> /known_hosts
RUN echo "UserKnownHostsFile /known_hosts" >> /etc/ssh/ssh_config

RUN git clone https://github.com/rbenv/ruby-build.git
RUN PREFIX=/usr/local ./ruby-build/install.sh
ENV RUBY_CONFIGURE_OPTS="ac_cv_func_isnan=yes ac_cv_func_isinf=yes"
RUN ruby-build "$RUBY_VERSION" /usr/

RUN gem update --system "$RUBYGEMS_VERSION"
RUN gem install bundler --version "$BUNDLER_VERSION" --force

ENV GEM_HOME /usr/local/bundle
ENV BUNDLE_PATH="$GEM_HOME" \
	BUNDLE_BIN="$GEM_HOME/bin" \
	BUNDLE_SILENCE_ROOT_WARNING=1 \
	BUNDLE_APP_CONFIG="$GEM_HOME"
ENV PATH $BUNDLE_BIN:$PATH
RUN mkdir -p "$GEM_HOME" "$BUNDLE_BIN" \
	&& chmod 777 "$GEM_HOME" "$BUNDLE_BIN"
COPY in.rb /opt/resource/in
COPY check.rb /opt/resource/check
COPY out.rb /opt/resource/out
COPY Gemfile /opt/resource/Gemfile
RUN cd /opt/resource/ && bundle install
