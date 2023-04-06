FROM gcr.io/graphite-docker-images/downstream-builder:latest

RUN apt-get install -y python3-openssl rsync
RUN curl -L https://github.com/pyenv/pyenv-installer/raw/master/bin/pyenv-installer | bash
ENV PATH="/root/.pyenv/bin:${PATH}"
RUN eval "$(pyenv init -)"
RUN eval "$(pyenv virtualenv-init -)"
RUN pyenv install 3.6.8
RUN pyenv rehash
ENV PATH="/root/.pyenv/shims:${PATH}"
RUN pyenv global 3.6.8