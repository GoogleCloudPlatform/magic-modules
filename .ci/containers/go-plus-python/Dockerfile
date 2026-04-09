FROM gcr.io/graphite-docker-images/go-plus:latest

# Install python & python libraries.
RUN apt-get update
RUN apt-get install -y git
RUN apt-get install -y rsync
RUN apt-get install -y build-essential libbz2-dev libssl-dev libreadline-dev \
                        libffi-dev libsqlite3-dev tk-dev
RUN apt-get install -y libpng-dev libfreetype6-dev
RUN apt-get install -y make build-essential libssl-dev zlib1g-dev libbz2-dev \
    libreadline-dev libsqlite3-dev llvm libncurses5-dev libncursesw5-dev \
    xz-utils tk-dev libffi-dev liblzma-dev python3-openssl
ENV LC_ALL=C.UTF-8
ENV LANG=C.UTF-8
