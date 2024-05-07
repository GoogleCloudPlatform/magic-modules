---
subcategory: "Cloud Composer"
description: |-
  Kubernetes secret in Composer Environment workloads.
---

# google\_composer\_user\_workloads\_secret

Provides access to Kubernetes secret configuration for a given project, region and Composer Environment.

## Example Usage

```hcl
resource "google_composer_environment" "example" {
    name = "example-environment"
    config{
        software_config {
            image_version = "composer-3-airflow-2"
        }
    }
}

resource "google_composer_user_workloads_secret" "example" {
    environment = google_composer_environment.example.name
    name = "example-secret"
    data = {
        username: base64encode("username"),
        password: base64encode("password"),
    }
}

data "google_composer_user_workloads_secret" "example" {
    environment = google_composer_environment.example.name
    name = resource.google_composer_user_workloads_secret.example.name
}

output "debug" {
    value = data.google_composer_user_workloads_secret.example
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Secret.

* `environment` - (Required) Environment where the secret is stored.

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

* `region` - (Optional) The location or Compute Engine region of the environment.

## Attributes Reference

The following attributes are exported:

* `id` - An identifier for the resource in format `projects/{{project}}/locations/{{region}}/environments/{{environment}}/userWorkloadsSecrets/{{name}}`

