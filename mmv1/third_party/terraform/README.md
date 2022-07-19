# Handwritten

## Overview

The Google providers for Terraform have a large number of handwritten go files, primarily for resources written before Magic Modules was used with them. Most handwritten files are expected to stay handwritten indefinitely, although conversion to a generator may be possible for a limited subset of them.

Handwritten resources like [google_container_cluster](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/container_cluster) can be identified if they have source code present under the [mmv1/third_party/terraform/resources](./resources) folder or by the absence of the `AUTO GENERATED CODE` header in their [Go source](https://github.com/hashicorp/terraform-provider-google/blob/main/google/resource_container_cluster.go) in the downstream repositories. Handwritten datasources should be under the [mmv1/third_party/terraform/data_sources](./data_sources) folder, tests under the [mmv1/third_party/terraform/tests](./tests) folder and web documentaion under the [mmv1/third_party/terraform/website](./website) folder.

## Table of Contents
- [Contributing](#contributing)
  - [Shared Concepts](#shared-concepts)
	- [Resource](#resource)
	- [Datasource](#datasource)
	- [IAM Resources](#iam-resource)
	- [Test](#test)
	- [Documentation](#documentation)
	- [Beta Feature](#beta-feature)
		- [Add or update a beta future](#add-or-update-a-beta-feature)
		- [Test that use a beta feature](#test-that-use-a-beta-feature)
		- [Promote a beta feature](#promote-a-beta-feature)


## Contributing

We're glad to accept contributions to handwritten resources. Tutorials and guidance on making changes are available below.

### Shared Concepts

This section will serve as a point of reference for some shared concepts that
all handwritten files share. It's meant to be an introduction to our serialization
strategy and overview.

#### Serialization strategy
The go files within the directory files are copied literally to their respective providers.
Our serialization methodology may seem complicated but for the case of handwritten resources its quite
simple. Editing the file will change its counterpart downstream.

#### go and go.erb
Within the third party library you'll notice `go` and `go.erb` files.
Go files are native golang code while go.erb pass through ruby before
being serialized. The reason `go.erb` files exist are to protect certain
properties or fields from entering the `ga` provider. Thus you'll often see
lines like `<% unless version == 'ga' -%>` within the file. These blocks
will omit the enclosure from being output to the GA provider. In the
rare case where you are promoting all fields to `ga` and these blocks
are no longer needed you can remove the `.erb` extension.

#### Create, Read, Update, Delete
As far as terraform schema is concerned these are the functions we
need to provide for terraform to be able to provision and delete
resources. In editing any fields you'll likely be adding functionality to
these functions or implementing them wholesale.


#### Expanders and Flatteners
Expanders and flatteners are concepts created to simplify common patterns
and add conformity/code consistency. Essentially expanders are functions
used to segregate some translation from terraform representation to api representation.
We will use these to encapsulate this translation for blocks and/or complicated fields.
This allows our code to be concise and functionality to be readable and easily
apparent by separating these into their own functions. While expanders are used
for terraform to api, flatteners do just the opposite. Converting api to terraform.

Thus
* expanders - helper functions used for translating tf -> api representation
* flatteners - helper functions used for translating api -> tf representation

### Resource
While we no longer accept new handwritten resources except in rare cases. Understanding
how to edit and add to existing resources may be important for implementing new fields
or changing existing behavior.

To edit an existing resource to add a field there are four steps you'll go through.
1. Add the new field to the schema
1. Implement the respective flattener and/or expander for the new field
1. Add a testcase for the field or extend an existing one
1. Add documentation for the field to the respective markdown file

##### Adding a new field to the schema
To add a new field you will have to compare an existing resource
to it's respective rest api documentation. Dependant on how the api implements
the field we will in almost all cases mirror the structure. For example if there is
an `enabled` field nested under a ``IdentityServiceConfig`` block we will mirror
this within the schema.

Thus the block for terraform to utilize this field would then be
```terraform
resource "x" "y" {
  `identity_service_config{
    enabled = true
  }`
}
```

You might think it convoluted to provide such a structure. Why not
simply provide a single `enable_identity_service_config`. One constant
has echoed through our mind as terraform developers through the years.
Api's are ever evolving. Mirroring the api gives us the best chance to stay
in step with that evolution. Therefore if `IdentityServiceConfig` is extended with new
parameters in the future we can cleanly encapsulate those into the existing block(s).

As far as providing the field itself, it's fairly straightforward. Mirror the field
from the api and look to the other fields and the schema type in the SDK to see
what's available and how to structure it. For the documentation, copying the documentation
from the rest api will be the usual practice.

If you are adding a field that is an ENUM from the api standpoint its best practice
to provide it as as a string to the provider. This field will likely have values
added to it by the api and this future proofs our provider to support new values without
haven't to make new additions. There will be rare exceptions, but generally its a good
practice.

#### Implement the respective flattener and/or expander for the new field
Once you've added the field to the schema you will implement the corresponding
expander/flattener. See [expanders and flatters](#expanders-and-flatteners) for
more context on what these fields are used for. Essentially we will be editing the
read, create, and update operations to parse the schema and call the api to make
the changes to the state of the resource. Following existing patterns to create
this operation will be the best way to implement this. As there are many unique ways
to implement a given field we won't get into specifics.

#### Add a testcase for the field or extend an existing one
Once your field has been implemented, go to the corresponding test file for
your resource and extend it. If your field is updatable it's good practice to
have a two step apply to ensure that the field *can* be updated. You'll notice
a lot of our tests have a import state verify directly after apply. These
steps are important as they will essentially attempt to import the resource
you just provisioned and *verify* that the field values are consistent with the
applied state. Please test all fields you've added to the provider. It's important
for us to ensure all fields are usable and workable.

#### Add documentation for the field to the respective markdown file
See [Documentation](#documentation) for more information. Essentially you will
just be opening the corresponding markdown file and adding documentation, likely
copied from the rest api to the markdown file. Follow the existing patterns there-in.

### Datasource
Datasources are like terraform resources except they don't *create* anything.
They are simply read-only operations that will expose some sort of values needed
for subsequent resource operations. If you're adding a field to an existing
datasource, check the [Resource](#resource) section. Everything there will
be mostly consistent with the type of change you'll need to make. For adding
a new datasource there are 5 steps to doing so.

1. Create a new datasource declaration file and a corresponding test file
1. Add Schema and Read operation implementation
1. Add the datasource to the `provider.go.erb` index
1. Implement a test which will create and resources and read the corresponding
  datasource.
1. Add documentation.

#### Create a new datasource declaration file and a corresponding test file

#### Add Schema and Read operation implementation

#### Add the datasource to the `provider.go.erb` index

#### Implement a test

#### Add documentation

### IAM Resource

### Test

### Documentation

### Beta Feature

#### Add or update a beta feature

#### Test that use a beta feature

#### Promote a beta feature