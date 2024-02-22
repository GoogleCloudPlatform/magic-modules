---
page_title: project_from_id Function - terraform-provider-google
description: |-
  Returns the project within a provided resource id, self link, or OP style resource name.
---

# Function: project_from_id

Returns the project within a provided resource id, self link, or OP style resource name.

## Example Usage

```terraform
terraform {
	required_providers {
		google = {
			source = "hashicorp/google"
		}
	}
}

# Value is "my-project"
output "function_output" {
	value = provider::google::project_from_id("https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-c/instances/my-instance")
}
```

## Signature

```text
project_from_id(id string) string
```

## Arguments

1. `id` (String) A string of a resource's id, a resource's self link, or an OP style resource name. These are all valid values:

* `"projects/my-project/zones/us-central1-c/instances/my-instance"`
* `"https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-c/instances/my-instance"`
* `"//gkehub.googleapis.com/projects/my-project/locations/us-central1/memberships/my-membership"`
