<img src='mm_logo.png' align='right'>

# Magic Modules


## Overview

Magic Modules (MM) is a project to auto-generate management libraries (e.g.
plugins) for third party tools to operate Google Cloud Platform (GCP) products
natively.


## Motivation

Customers -- large and small -- use many different deployment tools to provision
and manage their operating environments. Companies like Puppet, Chef, Ansible,
Hashicorp produce open source tools to meet these needs and typically provide
enterprise versions and support.

Once customers invest in these tools they prefer to use a single management
system (and style). For this reason Google wants to provide integration for
these tools, so that our customers have an easier path to GCP versus abandoning
a tool they're already comfortable with. For example Puppet can provision and
manage VMs in Compute Engine, databases in Cloud SQL, Google Kubernetes Engine
clusters, and so on.

But given the breadth of products available in Google Cloud Platform, creating
modules by hand is quickly becomes inefficient and costly. It would involve too
much redundant engineering work and create too much code for a team to maintain.
There is also the issue of explosion of modules, if you consider for
comprehensive coverage: <code>
total\_libraries&nbsp;=&nbsp;total\_products&nbsp;\*&nbsp;total\_providers
</code>.
The more products we add to the portfolio, or increase the support to other
tools, doing this work by hand will not scale.

Enter Magic Modules.


## Philosophy

The philosophy behind Magic Modules is a strongly defined object and object
relationship schema. This schema can in turn be used to map the necessary API
calls to achieve the object state.

> Magic Modules heavily relies on *object convergence* as its target
> destination. That means that the user will specify the final state they would
> like their objects (and/or object relations to be) and the code generated by
> Magic Modules will do "whatever it takes" to make it so.

Note that there is not a required 1:1 mapping between objects and API calls to
GCP. This 1:N mapping is key to the simplicity of the interface exposed to
users. There are various wrappers out there that can envelope a Google API into
various languages, e.g. Ruby or Java, but the user needs to know all the
intricacies of the underlying GCP product API.

[Read more about piylosophy][philosophy]

## Current Coverage

### Supported Providers

- Puppet
- Chef

### Supported Products

All the necessary data for a product to compile is under the `products/` folder.
For simplicity (yet not a requirement) the folder name is the actual product
name you will use in the `-p` parameter later on.

To see the full product list visit the [`products`](products/) folder.


## (One Time) Setup

Magic Modules requires:

- git (so you can checkout the code)
- Ruby 2.0 or higher
- Bundler gem

### Downloading code

Depending on the product and provider combination the generated code is stored
on another Git repository. To make it easier to track them all these foreign Git
repositories are tracked as submodules of Magic Modules main repository.

To fetch all products last known good build, run:

    git submodule update --init

### Installing bundler

As you have Ruby installed all you need is Bundler. To install it, run:

    gem install bundler

### Dependencies install

Now that we have Bundler installed, go to the root folder where you checked out
the Magic Modules code and run:

    bundle install

We are now ready to compile a module!


## Compiling a module

Compiling a module is as easy as:

    bundle exec ruby compiler.rb -p <product> -e <provider> -o <folder>

For example, to compile Google Container Engine for Puppet, you invoke:

    bundle exec ruby compiler.rb -p products/container -e puppet -o build/puppet/container

And the generated code should be written to `build/puppet/container`


## Expected output

You point Magic Modules to an _empty folder_ and it will output a fully
operational module for the selected provider.

Magic modules will compile the schema and generate the following output:

1. The production-grade, user facing code in the target tool language and style
2. Documentation for all objects and interactions between objects
3. All tests that certifies the user facing code
4. All API layer data used during test runs
5. All required libraries for the operation


## Testing the generated code

Each provider has their own testing mechanism. Run the command in the product's
output folder `build/<provider>/<product>`, e.g. `build/puppet/compute`.

Platform | Tool  | Test command
---------|-------|--------------
Puppet   | rspec | bundle exec rspec
Chef     | rspec | bundle exec rspec

And refer to the module documentation on how to execute and debug the code or
tests further.


## Creating or updating a module

Please refer to [Developer Guide][developer] for details on how to create a new
product or update an existing one.


[developer]: DEVELOPER.md
[philosophy]: docs/philosophy.md
