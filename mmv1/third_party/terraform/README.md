# Handwritten

## Overview

The Google providers for Terraform have a large number of handwritten go files, primarily for resources written before Magic Modules was used with them. Most handwritten files are expected to stay handwritten indefinitely, although conversion to a generator may be possible for a limited subset of them.

Handwritten resources like [google_container_cluster](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/container_cluster) can be identified if they have source code present under the [mmv1/third_party/terraform/resources](./resources) folder or by the absence of the `AUTO GENERATED CODE` header in their [Go source](https://github.com/hashicorp/terraform-provider-google/blob/main/google/resource_container_cluster.go) in the downstream repositories. Handwritten datasources should be under the [mmv1/third_party/terraform/data_sources](./data_sources) folder, tests under the [mmv1/third_party/terraform/tests](./tests) folder and web documentaion under the [mmv1/third_party/terraform/website](./website) folder.

## Table of Contents
- [Contributing](#contributing)
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

### Resource

### Datasource

### IAM Resource

### Test

### Documentation

### Beta Feature

#### Add or update a beta feature

#### Test that use a beta feature

#### Promote a beta feature