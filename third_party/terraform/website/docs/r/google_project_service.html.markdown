---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_project_service"
sidebar_current: "docs-google-project-service-x"
description: |-
 Allows management of a single API service for a Google Cloud Platform project.
---

# google\_project\_service

Allows management of a single API service for an existing Google Cloud Platform project. 

For a list of services available, visit the
[API library page](https://console.cloud.google.com/apis/library) or run `gcloud services list`.

~> **Note:** Previous (pre-3.0.0) versions of the provider had a `google_project_services` resource, which was authoritative over the list of GCP services enabled on a project, so that services not explicitly set would be removed by Terraform. However, this was dangerous to use in practice. Services have dependencies that are automatically enabled alongside them and GCP will add dependencies to services out of band, enabling them. If a user ran Terraform after this, Terraform would disable the service and implicitly disable any service that relied on it.The `google_project_service` resource is a much better match for most users' intent, managing a single service at a time. Setting several `google_project_service` resources is an assertion that "these services are set on this project", while `google_project_services` was an assertion that "only these services are set on this project".

## Example Usage

```hcl
resource "google_project_service" "project" {
  project = "your-project-id"
  service = "iam.googleapis.com"

  disable_dependent_services = true
}
```

## Argument Reference

The following arguments are supported:

* `service` - (Required) The service to enable.

* `project` - (Optional) The project ID. If not provided, the provider project is used.

* `disable_dependent_services` - (Optional) If `true`, services that are enabled and which depend on this service should also be disabled when this service is destroyed.
If `false` or unset, an error will be generated if any enabled services depend on this service when destroying it.

* `disable_on_destroy` - (Optional) If true, disable the service when the terraform resource is destroyed.  Defaults to true.  May be useful in the event that a project is long-lived but the infrastructure running in that project changes frequently.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `{{project}}/{{service}}`

## Import

Project services can be imported using the `project_id` and `service`, e.g.

```
$ terraform import google_project_service.my_project your-project-id/iam.googleapis.com
```

Note that unlike other resources that fail if they already exist, `terraform apply` can be successfully used to verify already enabled services. This means that when importing existing resources into Terraform, you can either import the `google_project_service` resources or treat them as new infrastructure and run `terraform apply` to add them to state.
