---
layout: "google"
page_title: "Terraform Google Provider 3.0.0 Upgrade Guide"
sidebar_current: "docs-google-provider-version-3-upgrade"
description: |-
  Terraform Google Provider 3.0.0 Upgrade Guide
---

# Terraform Google Provider 3.0.0 Upgrade Guide

The `3.0.0` release of the Google provider for Terraform is a major version and
includes some changes that you will need to consider when upgrading. This guide
is intended to help with that process and focuses only on the changes necessary
to upgrade from the final `2.X` series release to `3.0.0`.

Most of the changes outlined in this guide have been previously marked as
deprecated in the Terraform `plan`/`apply` output throughout previous provider
releases, up to and including the final `2.X` series release. These changes,
such as deprecation notices, can always be found in the CHANGELOG of the
affected providers. [google](https://github.com/terraform-providers/terraform-provider-google/blob/master/CHANGELOG.md)
[google-beta](https://github.com/terraform-providers/terraform-provider-google-beta/blob/master/CHANGELOG.md)

## I accidentally upgraded to 3.0.0, how do I downgrade to `2.X`?

If you've inadvertently upgraded to `3.0.0`, first see the
[Provider Version Configuration Guide](#provider-version-configuration) to lock
your provider version; if you've constrained the provider to a lower version
such as shown in the previous version example in that guide, Terraform will pull
in a `2.X` series release on `terraform init`.

If you've only ran `terraform init` or `terraform plan`, your state will not
have been modified and downgrading your provider is sufficient.

If you've ran `terraform refresh` or `terraform apply`, Terraform may have made
state changes in the meantime.

* If you're using a local state, or a remote state backend that does not support
versioning, `terraform refresh` with a downgraded provider is likely sufficient
to revert your state. The Google provider generally refreshes most state
information from the API, and the properties necessary to do so have been left
unchanged.

* If you're using a remote state backend that supports versioning such as
[Google Cloud Storage](https://www.terraform.io/docs/backends/types/gcs.html),
you can revert the Terraform state file to a previous version. If you do
so and Terraform had created resources as part of a `terraform apply` in the
meantime, you'll need to either delete them by hand or `terraform import` them
so Terraform knows to manage them.

## Upgrade Topics

<!-- TOC depthFrom:2 depthTo:2 -->

- [Provider Version Configuration](#provider-version-configuration)
- [Data Source: `google_container_engine_versions`](#data-source-google_container_engine_versions)
- [Resource: `google_cloudiot_registry`](#resource-google_cloudiot_registry)
- [Resource: `google_compute_forwarding_rule`](#resource-google_compute_forwarding_rule)
- [Resource: `google_compute_network`](#resource-google_compute_network)
- [Resource: `google_compute_network_peering`](#resource-google_compute_network_peering)
- [Resource: `google_compute_region_instance_group_manager`](#resource-google_compute_region_instance_group_manager)
- [Resource: `google_container_cluster`](#resource-google_container_cluster)
- [Resource: `google_container_node_pool`](#resource-google_container_node_pool)
- [Resource: `google_monitoring_alert_policy`](#resource-google_monitoring_alert_policy)
- [Resource: `google_monitoring_uptime_check_config`](#resource-google_monitoring_uptime_check_config)
- [Resource: `google_project_services`](#resource-google_project_services)
- [Resource: `google_storage_bucket`](#resource-google_storage_bucket)

<!-- /TOC -->

## Provider Version Configuration

-> Before upgrading to version 3.0.0, it is recommended to upgrade to the most
recent `2.X` series release of the provider and ensure that your environment
successfully runs [`terraform plan`](https://www.terraform.io/docs/commands/plan.html)
without unexpected changes or deprecation notices.

It is recommended to use [version constraints](https://www.terraform.io/docs/configuration/providers.html#provider-versions)
when configuring Terraform providers. If you are following that recommendation,
update the version constraints in your Terraform configuration and run
[`terraform init`](https://www.terraform.io/docs/commands/init.html) to download
the new version.

If you aren't using version constraints, you can use `terraform init -upgrade`
in order to upgrade your provider to the latest released version.

For example, given this previous configuration:

```hcl
provider "google" {
  # ... other configuration ...

  version = "~> 2.17.0"
}
```

An updated configuration:

```hcl
provider "google" {
  # ... other configuration ...

  version = "~> 3.0.0"
}
```

## Data Source: `google_container_engine_versions`

### `region` and `zone` are now removed

Use `location` instead.

## Resource: `google_cloudiot_registry`

### `event_notification_config` is now removed

`event_notification_config` has been removed in favor of
`event_notification_configs` (plural). Please switch to using the plural field.

## Resource: `google_compute_forwarding_rule`

### `ip_version` is now removed

`ip_version` is not used for regional forwarding rules.

## Resource: `google_compute_network`

### `ipv4_range` is now removed

Legacy Networks are deprecated and you will no longer be able to create them
using this field from Feb 1, 2020 onwards.

## Resource: `google_compute_network_peering`

### `auto_create_routes` is now removed

`auto_create_routes` has been removed because it's redundant and not
user-configurable.

## Resource: `google_compute_region_instance_group_manager`

### `update_strategy` no longer has any effect and is removed

With `rolling_update_policy` removed, `update_strategy` has no effect anymore.
Before updating, remove it from your config.

## Resource: `google_container_cluster`

### `zone`, `region` and `additional_zones` are now removed

`zone` and `region` have been removed in favor of `location` and
`additional_zones` has been removed in favor of `node_locations`

## Resource: `google_container_node_pool`

### `zone` and `region` are now removed

`zone` and `region` have been removed in favor of `location`

## Resource: `google_monitoring_alert_policy`

### `labels` is now removed

`labels` is removed as it was never used. See `user_labels` for the correct field.

## Resource: `google_monitoring_uptime_check_config`

### `is_internal` and `internal_checker` are now removed

`is_internal` and `internal_checker` never worked, and are now removed.

## Resource: `google_project_services`

### `google_project_services` has been removed from the provider

The `google_project_services` resource was authoritative over the list of GCP
services enabled on a project, so that services not explicitly set would be
removed by Terraform.

However, this was dangerous to use in practice. Services have dependencies that
are automatically enabled alongside them and GCP will add dependencies to
services out of band, enabling them. If a user ran Terraform after this,
Terraform would disable the service- and implicitly disable any service that
relied on it.

The `google_project_service` resource is a much better match for most users'
intent, managing a single service at a time. Setting several
`google_project_service` resources is an assertion that "these services are set
on this project", while `google_project_services` was an assertion that "**only**
these services are set on this project".

Users should migrate to using `google_project_service` resources, or using the
[`"terraform-google-modules/project-factory/google//modules/project_services"`](https://registry.terraform.io/modules/terraform-google-modules/project-factory/google/3.3.0/submodules/project_services)
module for a similar interface to `google_project_services`.

#### Old Config

```hcl
resource "google_project_services" "project" {
  project            = "your-project-id"
  services           = ["iam.googleapis.com", "cloudresourcemanager.googleapis.com"]
  disable_on_destroy = false
}
```

#### New Config (module)

```hcl
module "project_services" {
  source  = "terraform-google-modules/project-factory/google//modules/project_services"
  version = "3.3.0"

  project_id    = "your-project-id"
  activate_apis =  [
    "iam.googleapis.com",
    "cloudresourcemanager.googleapis.com",
  ]

  disable_services_on_destroy = false
  disable_dependent_services  = false
}
```

#### New Config (google_project_service)

```hcl
resource "google_project_service" "project_iam" {
  project = "your-project-id"
  service = "iam.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "project_cloudresourcemanager" {
  project = "your-project-id"
  service = "cloudresourcemanager.googleapis.com"
  disable_on_destroy = false
}
```

## Resource: `google_storage_bucket`

### `is_live` is now removed

Please use `with_state` instead, as `is_live` is now removed.
