---
page_title: "Managing Labels and Annotations with Terraform"
description: |-
  Manage resource labels and annotations with Terraform Google Provider 5.0
---

# Resource Labels and Annotations

Labels and annotations are key-value pairs attached on Google cloud resources. Cloud labels are used for organizing resources, filtering resources, breaking down billing, and so on. Annotations are used to attach metadata to Kubernetes resources.

Not all of Google Cloud resources support labels and annotations. Please check the Terraform Google provider resource documentation to figure out if a given resource supports `labels` or `annotations` fields.

## Managing Labels with Terraform Google Provider 5.0

The new labels model introdued in Terraform Google Provider 5.0 will be applied to all of the resources with the top level `labels` field or the nested `labels` field inside the top level `metadata` field. Some labels fields are for child resources, so the new model will not be applied to labels fields for child resources.

There are now three label-related fields with the new model:

* The `labels` field is now non-authoritative and only manages the label keys defined in your configuration for the resource.
* The `terraform_labels` cannot be specified directly by the user. It merges the labels defined in the resource's configuration and the default labels configured in the provider block. If the same label key exists on both the resource level and provider level, the value on the resource will override the provider-level default.
* The output-only `effective_labels` will list all the labels present on the resource in GCP, including the labels configured through Terraform, the system, and other clients.

### Import and managing resource labels

#### Example configuration
```hcl
resource "google_bigquery_dataset" "dataset" {
  dataset_id = "example_dataset"

  labels = {
    key1 = "value1"
    key2 = "value2"
  }
}
```

In this example, after running `terraform import`, `labels` and `terraform_labels` are empty, and `effective_labels` has all of labels present on the resource in GCP, `key1` and `key2`.

To bring labels defined in the configuration under management, run `terraform apply`. After the configuration is applied, `labels` is managing the user defined labels `key1` and `key1`. `terraform_labels` also has the user defined labels`key1` and `key2`. `effective_labels` is not affected.

### Add resource labels under management to an existing resource

#### Resource configuration without labels
```hcl
resource "google_dataproc_cluster" "cluster" {
  name   = "tf-test-dproc-test-1"
  region = "us-central1"
}
```

In this example, after the configuration is applied, resource `google_dataproc_cluster.cluster` is created without the `labels` field. `terraform_labels` field is empty and `effective_labels` only has system labels as neither resource labels nor provider default labels are configured.


#### Resource configuration with labels added
```hcl
resource "google_dataproc_cluster" "cluster" {
  name   = "tf-test-dproc-test-1"
  region = "us-central1"

  # The labels block is added.
  labels = {
    key1 = "value1"
    key2 = "value2"
  }
}
```

After the configuration with labels block is applied, Terraform is managing `key1` and `key2` in the `labels` field. `terraform_labels` field also has labels `key1` and `key2`. `effective_labels` has all of labels present on the resource in GCP, including `key1`, `key2` and system labels.

### Removing a managed label
#### Example configuration
```hcl
resource "google_dataproc_cluster" "cluster" {
  name   = "tf-test-dproc-test-1"
  region = "us-central1"

  labels = {
    key1 = "value1"
    # key2 is removed
  }
}
```

Applying this configuration **after** the previous example, `key2` is removed from all of three label related fields, `labels`, `terraform_labels` and `effective_labels`. Other values are unaffected.

### Ignoring labels in resources

[ignore_changes](https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#ignore_changes) can be applied to `labels` field to ignore the changes of the user defined labels.

```hcl
resource "google_dataproc_cluster" "cluster" {
  name   = "tf-test-dproc-test-1"
  region = "us-central1"

  labels = {
    key1 = "value1"
    key2 = "value2"
  }

  lifecycle {
    ignore_changes = [
      labels,
    ]
  }
}
```

In this example, the `key1` and `key2` labels will be added to the resource `google_dataproc_cluster.cluster` on resource creation, however any changes to the labels block will be ignored:

~> **Note:** It is not recommended to apply `ignore_changes` to `terraform_labels` or `effective_labels`, as it may unintuitively affect the final API call.

### Managing Out-of-Band labels

If a label was added outside of Terraform, it will not be managed by Terraform. To bring it under management, the out-of-band label needs to be added to the `labels` field in the configuration.

#### Example configuration
```hcl
resource "google_dataproc_cluster" "cluster" {
  name   = "tf-test-dproc-test-1"
  region = "us-central1"

  labels = {
    key1 = "value1"
    # Add the out-of-band label into configuration
    out_of_band = "value2"
  }
}
```

After the configuration is applied, Terraform starts to manage the label key `out_of_band`.

### Applying provider default labels

```hcl
provider "google" {
  default_labels = {
    my_global_key = "one"
    my_default_key = "two"
  }
}

resource "google_dataproc_cluster" "cluster" {
  name   = "tf-test-dproc-test-1"
  region = "us-central1"

  labels = {
    key1           = "value1"
    # overrides provider-wide setting
    my_default_key = "four"
  }
}
```

In this example, after the configuration is applied, Terraform is managing `key1` and `my_default_key` in the `labels` field. `terraform_labels` field has `key1`, `my_default_key` and `my_global_key`. `effective_labels` has all of labels present on the resource in GCP, including `key1`, `my_default_key`, `my_global_key` and system labels.

## Managing Annotations with Terraform Google Provider 5.0

The new annotations model introdued in Terraform Google Provider 5.0 will be applied to all of the resources with the top level `annotations` field or the nested `annotations` field inside the top level `metadata` field. Some annotations fields are for child resources, so the new model will not be applied to annotations fields for child resources.

There are now two annotation-related fields with the new model:

* The `annotations` field is now non-authoritative and only manages the label keys defined in your configuration for the resource.
* The output-only `effective_annotations` will list all the annotations present on the resource in GCP, including the annotations configured through Terraform, the system, and other clients.

Managing the annotations is similar to manage labels, except that provider-level default annotations are not supported.
