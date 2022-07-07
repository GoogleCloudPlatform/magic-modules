# MMV1

## Overview

Mmv1 is a code generator that implements the Terraform Google Provider (TGP) resources from ruby scripts.

Mmv1-generated resources like [google_compute_address](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_address) can be identified by looking in their [Go source](https://github.com/hashicorp/terraform-provider-google/blob/main/google/resource_compute_address.go) for an `AUTO GENERATED CODE` header as well as a Type `MMv1`. Mmv1-generated resources should have source code present under their product folder, like [mmv1/products/compute](./products/compute) for [google_compute_address](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_address) resource.

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

We'd love to accept your contributions! Thanks for making the changes :) Here's some guidance on how to contribute to mmv1-genereated resources.

### Resource

### IAM Resource

### Test

### Documentation

### Beta Feature

#### Add or update a beta feature

#### Test that use a beta feature

#### Promote a beta feature