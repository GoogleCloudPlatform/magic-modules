---
page_title: "Terraform Google Provider 5.0.0 Upgrade Guide"
description: |-
  Terraform Google Provider 5.0.0 Upgrade Guide
---

# Terraform Google Provider 5.0.0 Upgrade Guide

The `5.0.0` release of the Google provider for Terraform is a major version and
includes some changes that you will need to consider when upgrading. This guide
is intended to help with that process and focuses only on the changes necessary
to upgrade from the final `4.X` series release to `5.0.0`.

Most of the changes outlined in this guide have been previously marked as
deprecated in the Terraform `plan`/`apply` output throughout previous provider
releases, up to and including the final `4.X` series release. These changes,
such as deprecation notices, can always be found in the CHANGELOG of the
affected providers. [google](https://github.com/hashicorp/terraform-provider-google/blob/main/CHANGELOG.md)
[google-beta](https://github.com/hashicorp/terraform-provider-google-beta/blob/main/CHANGELOG.md)

## I accidentally upgraded to 5.0.0, how do I downgrade to `4.X`?

If you've inadvertently upgraded to `5.0.0`, first see the
[Provider Version Configuration Guide](#provider-version-configuration) to lock
your provider version; if you've constrained the provider to a lower version
such as shown in the previous version example in that guide, Terraform will pull
in a `4.X` series release on `terraform init`.

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
[Google Cloud Storage](https://developer.hashicorp.com/terraform/language/settings/backends/gcs),
you can revert the Terraform state file to a previous version. If you do
so and Terraform had created resources as part of a `terraform apply` in the
meantime, you'll need to either delete them by hand or `terraform import` them
so Terraform knows to manage them.

## Provider Version Configuration

-> Before upgrading to version 5.0.0, it is recommended to upgrade to the most
recent `4.X` series release of the provider, make the changes noted in this guide,
and ensure that your environment successfully runs
[`terraform plan`](https://developer.hashicorp.com/terraform/cli/commands/plan)
without unexpected changes or deprecation notices.

It is recommended to use [version constraints](https://developer.hashicorp.com/terraform/language/providers/requirements#requiring-providers)
when configuring Terraform providers. If you are following that recommendation,
update the version constraints in your Terraform configuration and run
[`terraform init`](https://developer.hashicorp.com/terraform/cli/commands/init) to download
the new version.

If you aren't using version constraints, you can use `terraform init -upgrade`
in order to upgrade your provider to the latest released version.

For example, given this previous configuration:

```hcl
terraform {
  required_providers {
    google = {
      version = "~> 4.70.0"
    }
  }
}
```

An updated configuration:

```hcl
terraform {
  required_providers {
    google = {
      version = "~> 5.0.0"
    }
  }
}
```

## Provider-level Labels Rework

The labels and annotations are key-value pairs attached on Google cloud resources. Cloud labels are used for organizing resources, filtering resources, breaking down billing, and so on. Annotations are used to attach metadata to Kubernetes resources.

Not all of Google cloud resources support labels and annotations. Please check the Terraform Google provider resource documentation to figure out if the resource supports the `labels` and `annotations` fields.

### Provider default labels

Default labels configured on the provider through the new `default_labels` field are now supported. The default labels configured on the provider will be applied to all of the resources with the top level `labels` field or the nested `labels` field inside the `metadata` field.

Provider-level default annotations are not supported.

### Resource labels

The new labels model will be applied to all of the resources with the top level `labels` field or the nested `labels` field inside the `metadata` field. Some labels fields are for child resources, so the new model will not be applied to the `labels` field for child resources.

There are now three label-related fields with the new model:

* The `labels` field will be non-authoritative and only manage the labels defined by the users on the resource through Terraform.
* The output-only `effective_labels` will list all of labels present on the resource in GCP, including the labels configured through Terraform, other clients and services.
* The output-only `terraform_labels` will merge the labels defined by the users on the resource through Terraform and the default labels configured on the provider. If the same label exists on both the resource labels and provider default labels, the label on the resource will override the provider label.

After upgrading to `5.0.0`, and then running `terraform refresh` or `terraform apply`, these three fields should show in the state file of the resources with a self-applying `labels` field.

### Resource annotations

The new annotations model is similar to the new labels model. 

There will be two annotation-related fields on resources with resource `annotations` field, the `annotations` and the output-only `effective_annotations` fields.

## Provider

### Provider-level change example header

Description of the change and how users should adjust their configuration (if needed).

## Datasources

## Datasource: `google_product_datasource`

### Datasource-level change example header

Description of the change and how users should adjust their configuration (if needed).

## Resources

## Resource: `google_product_resource`

### Resource-level change example header

Description of the change and how users should adjust their configuration (if needed).

## Resource: `google_firebase_web_app`

### `deletion_policy` now defaults to `DELETE`

Previously, `google_firebase_web_app` deletions default to `ABANDON`, which means to only stop tracking the WebApp in Terraform. The actual app is not deleted from the Firebase project. If you are relying on this behavior, set `deletion_policy` to `ABANDON` explicitly in the new version.
