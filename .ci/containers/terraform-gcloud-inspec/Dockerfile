FROM gcr.io/magic-modules/go-ruby-python:1.11.5-2.6.0-2.7-v6

RUN apt-get install unzip
RUN curl https://releases.hashicorp.com/terraform/0.12.16/terraform_0.12.16_linux_amd64.zip > terraform_0.12.16_linux_amd64.zip
RUN unzip terraform_0.12.16_linux_amd64.zip -d /usr/bin
# Install google cloud sdk
RUN echo "deb http://packages.cloud.google.com/apt cloud-sdk-stretch main" >> /etc/apt/sources.list.d/google-cloud-sdk.list
RUN curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
RUN apt-get update && apt-get install google-cloud-sdk -y

ADD Gemfile Gemfile
RUN bundle install
