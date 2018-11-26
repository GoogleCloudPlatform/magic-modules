FROM nmckinley/go-ruby-python:1.11-2.5.1-2.7

RUN apt-get install unzip
RUN curl https://releases.hashicorp.com/terraform/0.11.10/terraform_0.11.10_linux_amd64.zip > terraform_0.11.10_linux_amd64.zip
RUN unzip terraform_0.11.10_linux_amd64.zip -d /usr/bin
# Install google cloud sdk
RUN echo "deb http://packages.cloud.google.com/apt cloud-sdk-stretch main" >> /etc/apt/sources.list.d/google-cloud-sdk.list
RUN curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
RUN apt-get update && apt-get install google-cloud-sdk -y

# Install inspec
RUN gem install inspec -v '~> 3.0.0'