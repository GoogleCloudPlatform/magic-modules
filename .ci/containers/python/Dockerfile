from python:2.7-stretch

run pip install pygithub
run pip install absl-py
run pip install autopep8
run pip install beautifulsoup4 mistune

# Set up Github SSH cloning.
RUN ssh-keyscan github.com >> /known_hosts
RUN echo "UserKnownHostsFile /known_hosts" >> /etc/ssh/ssh_config
