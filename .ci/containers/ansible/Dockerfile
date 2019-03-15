FROM gcr.io/magic-modules/go-ruby-python:1.11.5-2.6.0-2.7

# Install python & python libraries.
RUN apt-get update
RUN apt-get install -y python python-pip
RUN pip install beautifulsoup4 mistune

# Install man for hacking/env-setup
RUN apt-get install -y man

# Install python dependencies for ansible
RUN pip install google-auth requests tox
