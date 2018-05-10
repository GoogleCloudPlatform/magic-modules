# Magic Modules Swiss Knife (SK)

A docker container with all the Magic Modules providers installed. 

## Requirements

* docker

## Providers

* Puppet
* Ansible *(Not yet supported)* 
* Chef *(Not yet supported)*
* Terraform *(Not yet supported)*

## Usage

### Step 1: Create credentials

1. In the Google Cloud console, go to IAM & admin > Service accounts
1. Create a new service account and download the json credential file.
1. Grant the required permissions to the service account.

### Step 2: Launch the container

In a terminal on your development machine, run:

```sh
# Specify the path to your credential file and magic-modules checkout.
export SK_CRED=~/path/to/your/cred/file
export SK_MM_PATH=~/path/to/magic-modules

# Generate the Magic Module for all providers
cd $SK_MM_PATH
bundle exec rake compile 

# Run the Swiss Knife container and start a bash session
docker run -i -t -v $SK_MM_PATH:/opt/magic-modules -v $SK_CRED:/etc/creds.json gcr.io/magic-modules/swiss-knife
```

### Step 3: Test your changes

#### Puppet

Inside the Swiss Knife container, run:

```sh
puppet apply /opt/magic-modules/build/puppet/compute/examples/delete_target_https_proxy.pp
puppet apply /opt/magic-modules/build/puppet/compute/examples/target_https_proxy.pp
```

#### Chef

Not yet supported

#### Terraform

Not yet supported

#### Ansible

Not yet supported

### Step 4: Rinse, repeat

Your `magic-modules` directory on your development machine is shared with the
Swiss Knife docker container.

This means you can simply make changes and call `bundle exec rake compile` and
the Swiss Knife will pick up the latest change.

## Developing the Swiss Knife

If you wish to add or improve support for any providers, 

### Building a new image

```sh
cd `/path/to/magic-modules/tools/swissknife
docker build -t swiss-knife .
```

### Publishing a new image

This requires extra permissions to write to the storage bucket.

```sh
# Replace <VERSION> with the version you want. e.g.: 0.1
docker tag swiss-knife gcr.io/magic-modules/swiss-knife:<VERSION>
gcloud docker -- push gcr.io/magic-modules/swiss-knife:<VERSION>
```
