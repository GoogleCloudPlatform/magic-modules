# Handwritten

## Overview

The Google providers for Terraform have a large number of handwritten go files, written before Magic Modules was used with them. While conversion is ongoing, many resources are still managed by hand.

Handwritten resources like [google_container_cluster](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/container_cluster) can be identified if they have source code present under the [mmv1/third_party/terraform/resources](./resources) folder or by the absence of the `AUTO GENERATED CODE` header in their [Go source](https://github.com/hashicorp/terraform-provider-google/blob/main/google/resource_container_cluster.go) in the downstream repositories.

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

We'd love to accept your contributions! Thanks for making the changes :) Here's some guidance on how to contribute to handwritten resources.

### Resource

### Datasource

### IAM Resource

### Test

### Documentation

### Beta Feature

#### Add or update a beta feature

#### Test that use a beta feature

#### Promote a beta feature