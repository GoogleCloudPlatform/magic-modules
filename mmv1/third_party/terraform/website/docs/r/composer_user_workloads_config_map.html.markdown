---
subcategory: "Cloud Composer"
description: |-
  User workloads ConfigMap used by Airflow tasks that run with Kubernetes Executor or KubernetesPodOperator.
---

# google\_composer\_user\_workloads\_config\_map

~> **Warning:** These resources are in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.

User workloads ConfigMap used by Airflow tasks that run with Kubernetes Executor or KubernetesPodOperator. 
Intended for Composer 3 Environments.

## Example Usage

```hcl
resource "google_composer_environment" "example" {
  name              = "example-environment"
  project           = "example-project"
  config {
    software_config {
      image_version = "example-image-version"
    }
  }
}

resource "google_composer_user_workloads_config_map" "example" {
  name = "example-config-map"
  project = "example-project"
  environment = google_composer_environment.example.name
  data = {
    api_host:  "apihost:443"
    db_host:   "dbhost:5432"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  Name of the Kubernetes ConfigMap.

* `region` -
  (Optional)
  The location or Compute Engine region for the environment.

* `project` -
  (Optional)
  The ID of the project in which the resource belongs.
  If it is not provided, the provider project is used.

* `environment` -
  Environment where the Kubernetes ConfigMap will be stored and used.

* `data` -
  (Optional)
  The "data" field of Kubernetes ConfigMap, organized in key-value pairs.
  For details see: https://kubernetes.io/docs/concepts/configuration/configmap/



## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{region}}/environments/{{environment}}/userWorkloadsConfigMaps/{{name}}`

## Import

ConfigMap can be imported using any of these accepted formats:

* `projects/{{project}}/locations/{{region}}/environments/{{environment}}/userWorkloadsConfigMaps/{{name}}`
* `{{project}}/{{region}}/{{environment}}/{{name}}`
* `{{name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import User Workloads ConfigMap using one of the formats above. For example:

```tf
import {
  id = "projects/{{project}}/locations/{{region}}/environments/{{environment}}/userWorkloadsConfigMaps/{{name}}"
  to = google_composer_user_workloads_config_map.example
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), Environment can be imported using one of the formats above. For example:

```
$ terraform import google_composer_user_workloads_config_map.example projects/{{project}}/locations/{{region}}/environments/{{environment}}/userWorkloadsConfigMaps/{{name}}
$ terraform import google_composer_user_workloads_config_map.example {{project}}/{{region}}/{{environment}}/{{name}}
$ terraform import google_composer_user_workloads_config_map.example {{name}}
```
