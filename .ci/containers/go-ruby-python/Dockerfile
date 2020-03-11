FROM gcr.io/magic-modules/go-ruby:1.13.8-2.6.0-v2

# Install python & python libraries.
RUN apt-get update
RUN apt-get install -y git
RUN apt-get install -y rsync
RUN apt-get install -y build-essential libbz2-dev libssl-dev libreadline-dev \
                        libffi-dev libsqlite3-dev tk-dev
RUN apt-get install -y libpng-dev libfreetype6-dev
RUN apt-get install -y make build-essential libssl-dev zlib1g-dev libbz2-dev \
    libreadline-dev libsqlite3-dev wget curl llvm libncurses5-dev libncursesw5-dev \
    xz-utils tk-dev libffi-dev liblzma-dev python-openssl
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
