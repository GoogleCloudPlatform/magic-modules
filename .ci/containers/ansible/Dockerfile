FROM gcr.io/magic-modules/go-ruby-python:1.11.5-2.6.0-2.7

# Install python & python libraries.
RUN apt-get update
RUN apt-get install -y python python-pip
RUN pip install beautifulsoup4 mistune

# Install man for hacking/env-setup
RUN apt-get install -y man locales

# Install python dependencies for ansible
RUN pip install google-auth requests tox

# Make sure locales set properly.
# This ensures that Ansible can read in
# config files properly (it's an ascii issue)
RUN echo "LC_ALL=en_US.UTF-8" >> /etc/environment
RUN echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen
RUN echo "LANG=en_US.UTF-8" > /etc/locale.conf
RUN locale-gen en_US.UTF-8
