# MMv1

## Overview

MMv1 is a Ruby-based code generator that implements Terraform Provider Google (TPG) resources from YAML specification files.

MMv1-generated resources like [google_compute_address](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_address) can be identified by looking in their [Go source](https://github.com/hashicorp/terraform-provider-google/blob/main/google/resource_compute_address.go) for an `AUTO GENERATED CODE` header as well as a Type `MMv1`. MMv1-generated resources should have source code present under their product folders, like [mmv1/products/compute](./products/compute) for the [google_compute_address](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_address) resource.

## Table of Contents
- [Contributing](#contributing)
	- [Resource](#resource)
	- [IAM Resources](#iam-resource)
	- [Test](#test)
	- [Documentation](#documentation)
	- [Beta Feature](#beta-feature)
		- [Add or update a beta future](#add-or-update-a-beta-feature)
		- [Test that use a beta feature](#test-that-use-a-beta-feature)
		- [Promote a beta feature](#promote-a-beta-feature)

## Contributing

We're glad to accept contributions to MMv1-generated resources. Tutorials and guidance on making changes are available below.

### Resource

### IAM Resource

### Test

### Documentation

### Beta Feature

#### Add or update a beta feature

#### Test that use a beta feature

#### Promote a beta feature