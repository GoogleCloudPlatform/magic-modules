---
subcategory: "Data Lineage"
description: |-
  Static lineage represented in OpenLineage format.
---

# google_data_lineage_open_lineage_job

Static lineage represented in OpenLineage format.

Defines lineage between datasets in OpenLineage format and publishes it to Knowledge Catalog.


To get more information about OpenLineageJob, see:

* [API documentation](https://docs.cloud.google.com/dataplex/docs/reference/data-lineage/rpc/google.cloud.datacatalog.lineage.v1)
* How-to Guides
    * [Official Documentation](https://docs.cloud.google.com/dataplex/docs/open-lineage)

<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_image=gcr.io%2Fcloudshell-images%2Fcloudshell%3Alatest&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md&cloudshell_working_dir=data_lineage_open_lineage_job_simple&open_in_editor=main.tf" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>
## Example Usage - Data Lineage Open Lineage Job Simple


```hcl
resource "google_dataplex_openlineage_job" "simple" {
  namespace   = "example_simple_namespace"
  name        = "example_simple_name"
  description = "Nightly ETL from raw to curated"

  inputs {
    namespace = "gs://example-bucket/"
    name      = "warehouse/raw_dataset_simple/source_table_1"
  }

  outputs {
    namespace = "gs://example-bucket/"
    name      = "warehouse/target_simple/target_table_1"
  }
}
```
<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_image=gcr.io%2Fcloudshell-images%2Fcloudshell%3Alatest&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md&cloudshell_working_dir=data_lineage_open_lineage_job_with_facets&open_in_editor=main.tf" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>
## Example Usage - Data Lineage Open Lineage Job With Facets


```hcl
resource "google_dataplex_openlineage_job" "with_facets" {
  namespace   = "example_with_facets_namespace"
  name        = "example_with_facets_name"
  description = "Nightly ETL from raw to curated"

  ownership {
    owners {
      name = "team:data-engineering"
      type = "MAINTAINER"
    }
  }

  inputs {
    namespace = "gs://example-bucket/"
    name      = "warehouse/raw_dataset_with_facets/source_table_1"

    symlinks {
      identifier {
        namespace = "bigquery"
        name      = "my-project-name.raw_dataset_with_facets".source_table_1"
        type      = "TABLE"
      }
    }

    catalog {
      framework = "bigquery"
      type      = "TABLE"
      name      = "my-project-name"
    }
  }

  outputs {
    namespace = "gs://example-bucket/"
    name      = "warehouse/target_with_facets/target_table_1"

  symlinks {
      identifier {
        namespace = "bigquery"
        name      = "my-project-name.target_dataset_with_facets.target_table_1"
        type      = "TABLE"
      }
    }

    catalog {
      framework = "bigquery"
      type      = "TABLE"
      name      = "my-project-name"
    }

    column_lineage {
      fields {
        name = "user_id"
        input_field {
          namespace = "gs://example-bucket/"
          name      = "warehouse/raw_dataset_with_facets/source_table_1"
          field     = "id"
          transformation {
            type    = "DIRECT"
            subtype = "IDENTITY"
          }
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:


* `namespace` -
  (Required)
  Namespace of the OpenLineage job.

* `name` -
  (Required)
  Name of the OpenLineage job.


* `description` -
  (Optional)
  Description of the OpenLineage job.

* `owner` -
  (Optional)
  The owner of the OpenLineage job.
  Structure is [documented below](#nested_owner).

* `input` -
  (Optional)
  Input datasets consumed by this job.
  Structure is [documented below](#nested_input).

* `output` -
  (Optional)
  Output datasets produced by this job.
  Structure is [documented below](#nested_output).

* `deletion_policy` - (Optional) Whether Terraform will be prevented from destroying the resource. Defaults to DELETE.
	When a 'terraform destroy' or 'terraform apply' would delete the resource,
	the command will fail if this field is set to "PREVENT" in Terraform state.
	When set to "ABANDON", the command will remove the resource from Terraform
	management without updating or deleting the resource in the API.
	When set to "DELETE", deleting the resource is allowed.


<a name="nested_owner"></a>The `owner` block supports:

* `name` -
  (Required)
  Owner name.

* `type` -
  (Required)
  Owner type.

<a name="nested_input"></a>The `input` block supports:

* `namespace` -
  (Required)
  Namespace of the dataset.

* `name` -
  (Required)
  Name of the dataset.

* `symlink` -
  (Optional)
  Symlink targets for the dataset.
  Structure is [documented below](#nested_input_symlink).

* `catalog` -
  (Optional)
  Catalog information for the dataset.
  Structure is [documented below](#nested_input_catalog).


<a name="nested_input_symlink"></a>The `symlink` block supports:

* `namespace` -
  (Required)
  Namespace of the symlink target.

* `name` -
  (Required)
  Name of the symlink target.

* `type` -
  (Required)
  Type of the symlink target.

<a name="nested_input_catalog"></a>The `catalog` block supports:

* `framework` -
  (Required)
  Catalog framework.

* `type` -
  (Required)
  Catalog entity type.

* `name` -
  (Required)
  Catalog entity name.

<a name="nested_output"></a>The `output` block supports:

* `namespace` -
  (Required)
  Namespace of the dataset.

* `name` -
  (Required)
  Name of the dataset.

* `symlink` -
  (Optional)
  Symlink targets for the dataset.
  Structure is [documented below](#nested_output_symlink).

* `catalog` -
  (Optional)
  Catalog information for the dataset.
  Structure is [documented below](#nested_output_catalog).

* `column_lineage` -
  (Optional)
  Column-level lineage information for the output dataset.
  Structure is [documented below](#nested_output_column_lineage).


<a name="nested_output_symlink"></a>The `symlink` block supports:

* `namespace` -
  (Required)
  Namespace of the symlink target.

* `name` -
  (Required)
  Name of the symlink target.

* `type` -
  (Required)
  Type of the symlink target.

<a name="nested_output_catalog"></a>The `catalog` block supports:

* `framework` -
  (Required)
  Catalog framework.

* `type` -
  (Required)
  Catalog entity type.

* `name` -
  (Required)
  Catalog entity name.

<a name="nested_output_column_lineage"></a>The `column_lineage` block supports:

* `field` -
  (Required)
  Field-level lineage mappings.
  Structure is [documented below](#nested_output_column_lineage_field).

* `dataset_input` -
  (Required)
  Input fields participating in output dataset lineage.
  Structure is [documented below](#nested_output_column_lineage_dataset_input).


<a name="nested_output_column_lineage_field"></a>The `field` block supports:

* `name` -
  (Required)
  Output field name.

* `input` -
  (Required)
  Input fields contributing to this output field.
  Structure is [documented below](#nested_output_column_lineage_field_input).


<a name="nested_output_column_lineage_field_input"></a>The `input` block supports:

* `namespace` -
  (Required)
  Namespace of the source dataset.

* `name` -
  (Required)
  Name of the source dataset.

* `field` -
  (Required)
  Source field name.

* `transformation` -
  (Optional)
  Transformations applied from source to output field.
  Structure is [documented below](#nested_output_column_lineage_field_input_transformation).


<a name="nested_output_column_lineage_field_input_transformation"></a>The `transformation` block supports:

* `type` -
  (Required)
  Transformation type.

* `subtype` -
  (Optional)
  Transformation subtype.

<a name="nested_output_column_lineage_dataset_input"></a>The `dataset_input` block supports:

* `namespace` -
  (Required)
  Namespace of the source dataset.

* `name` -
  (Required)
  Name of the source dataset.

* `field` -
  (Required)
  Source field name.

* `transformation` -
  (Optional)
  Transformations applied to fields from this input.
  Structure is [documented below](#nested_output_column_lineage_dataset_input_transformation).


<a name="nested_output_column_lineage_dataset_input_transformation"></a>The `transformation` block supports:

* `type` -
  (Required)
  Transformation type.

* `subtype` -
  (Optional)
  Transformation subtype.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `{{parent}}/locations/{{location}}/{{process}}`

* `knowledge_catalog` -
  Knowledge Catalog entities generated for this lineage job.
  Structure is [documented below](#nested_knowledge_catalog).


<a name="nested_knowledge_catalog"></a>The `knowledge_catalog` block contains:

* `process` -
  (Output)
  Knowledge Catalog process identifier.

* `run` -
  (Output)
  Knowledge Catalog run identifier.

## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.
