---
subcategory: "Dataplex"
description: |-
  Manages Data Lineage entities in Dataplex by using OpenLineage event API on every Create and Update, tracking job-level data lineage in Dataplex.
---

# google_dataplex_lineage_job

Emits an [OpenLineage](https://openlineage.io/) event to Dataplex Lineage API
[GCP Data Lineage API](https://docs.cloud.google.com/dataplex/docs/reference/data-lineage/rpc/google.cloud.datacatalog.lineage.v1)
(`ProcessOpenLineageRunEvent`) on every `create` and `update`.

A Dataplex **Process** entity is created on the first `apply` and reused on
subsequent applies; each `apply` creates a new **Run** and **LineageEvent**.

The implementation is based on the [open-lineage-base-resource]() module which provides a framework for defining a
resource generating static lineage OpenLineage events.
The module also contains objects for generating terraform schema capable of representing the OpenLineage job/dataset
entities in terraform.

- `Capability` interface
  - defines which OL facets are supported (in schema and generated events)
- `SchemaGenerator`
  - generates terraform schema for representing OpenLineage entities (jobs, datasets, facets) in terraform config
  - configured with `Capability`, if facets are not enabled, the schema is stubbed with `Optional` + `Computed`
    attributes to enable portability of configurations across OL providers
- EventBuilder - generates OpenLineage events based on the resource model and configured `Capability`

The resource enables the following OpenLineage facets

| Scope                      | Facet           | Terraform block  |
|----------------------------|-----------------|------------------|
| Job                        | `JobType`       | `job_type`       |
| Job                        | `Ownership`     | `ownership`      |
| Dataset (inputs & outputs) | `Symlinks`      | `symlinks`       |
| Dataset (inputs & outputs) | `Catalog`       | `catalog`        |
| Dataset (outputs only)     | `ColumnLineage` | `column_lineage` |

All other OL facet blocks (`documentation`, `source_code`, `source_code_location`,
`sql`, `tags`, `schema`, `data_source`, `storage`, …) are accepted by the schema
as silent no-ops so that a config shared with other OL providers is not rejected.

## Example Usage - Dataplex Lineage Job Basic

```hcl
resource "google_dataplex_lineage_job" "basic" {
  project   = "my-project-name"
  location  = "us"
  namespace = "my-pipeline-namespace"
  name      = "my-etl-job"
}
```

## Example Usage - Dataplex Lineage Job With Facets

```hcl
resource "google_dataplex_lineage_job" "with_facets" {
  project     = "my-project-name"
  location    = "us"
  namespace   = "my-pipeline-namespace"
  name        = "my-etl-job"
  description = "Nightly ETL from raw to curated"

  job_type {
    processing_type = "BATCH"
    integration     = "BYOL"
  }

  ownership {
    owners {
      name = "team:data-engineering"
      type = "MAINTAINER"
    }
  }

  inputs {
    namespace = "bigquery"
    name      = "my-project-name.raw_dataset.source_table"

    symlinks {
      identifier {
        namespace = "bigquery"
        name      = "my-project-name.raw_dataset.source_table"
        type      = "TABLE"
      }
    }

    catalog {
      framework = "bigquery"
      type      = "TABLE"
      name      = "my-project-name.raw_dataset.source_table"
    }
  }

  outputs {
    namespace = "bigquery"
    name      = "my-project-name.curated_dataset.output_table"

    catalog {
      framework = "bigquery"
      type      = "TABLE"
      name      = "my-project-name.curated_dataset.output_table"
    }

    column_lineage {
      fields {
        name = "user_id"
        input_field {
          namespace = "bigquery"
          name      = "my-project-name.raw_dataset.source_table"
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

### Required

* `location` - (Required, ForceNew) GCP region or multi-region where the Data
  Lineage API operates (e.g. `"us"`, `"eu"`, `"us-central1"`). Changing this
  forces a new resource to be created.

* `namespace` - (Required, ForceNew) OpenLineage job namespace. Changing this
  forces a new resource to be created.

* `name` - (Required, ForceNew) OpenLineage job name. Changing this forces a
  new resource to be created.

### Optional

* `project` - (Optional) The GCP project ID. Defaults to the project configured
  on the provider.

* `description` - (Optional) Human-readable description of the job.

---

#### `job_type` block (Optional)

Job type classification (`facets.JobType`). When this block is present all
three attributes are required.

* `processing_type` - (Required in block) `BATCH` or `STREAMING`.
* `integration` - (Required in block) Integration type, e.g. `SPARK`, `AIRFLOW`,
  `DBT`, `BYOL`.
* `job_type` - (Required in block) Job type, e.g. `QUERY`, `DAG`, `TASK`, `JOB`,
  `MODEL`.

---

#### `ownership` block (Optional)

Job owners (`facets.OwnershipJobFacet`). Requires at least one nested
`owners` block.

**`owners` block (Required in `ownership`)**

* `name` - (Required in block) Owner identifier, e.g. `team:data-engineering`.
* `type` - (Required in block) Owner type, e.g. `MAINTAINER`, `OWNER`, `STEWARD`.

---

#### `inputs` block (Optional, repeatable)

Input datasets consumed by this job.

* `namespace` - (Required) Dataset namespace.
* `name` - (Required) Dataset fully-qualified name.

Each `inputs` block may contain the following nested blocks:

**`symlinks` block (Optional)**

Alternate dataset identifiers (`facets.Symlinks`).

* `identifier` block (repeatable, at least 1 required):
  * `namespace` - (Required in block) Alternate namespace.
  * `name` - (Required in block) Alternate name.
  * `type` - (Required in block) e.g. `TABLE`, `VIEW`.

**`catalog` block (Optional)**

Catalog/metastore registration (`facets.Catalog`). When present, `framework`,
`type`, and `name` are required.

* `framework` - (Required in block) e.g. `bigquery`, `hive`, `iceberg`.
* `type` - (Required in block) Catalog type, e.g. `TABLE`, `VIEW`.
* `name` - (Required in block) Catalog-qualified name.
* `metadata_uri` - (Optional) e.g. `hive://localhost:9083`.
* `warehouse_uri` - (Optional) e.g. `hdfs://localhost/warehouse`.
* `source` - (Optional) Source system, e.g. `spark`.
* `catalog_properties` - (Optional) Additional catalog-specific properties as a
  `map(string)`.

---

#### `outputs` block (Optional, repeatable)

Output datasets produced by this job. Accepts the same `symlinks` and `catalog`
nested blocks as `inputs`, plus the following:

**`column_lineage` block (Optional)**

Column-level lineage for this output dataset (`facets.ColumnLineage`). Requires
at least one `fields` block.

* `fields` block (repeatable, at least 1 required):
  * `name` - (Required in block) Output column name.
  * `input_field` block (repeatable, at least 1 required):
    * `namespace` - (Required in block) Input dataset namespace.
    * `name` - (Required in block) Input dataset name.
    * `field` - (Required in block) Input column name.
    * `transformation` block (Optional, repeatable):
      * `type` - (Required in block) `DIRECT` or `INDIRECT`.
      * `subtype` - (Required in block) e.g. `IDENTITY`, `AGGREGATION`,
        `FILTER`.
      * `description` - (Optional) Human-readable transformation description.
      * `masking` - (Optional) `true` if this transformation masks/anonymises
        data. Defaults to `false`.

* `dataset` block (Optional, repeatable) — dataset-level lineage when the
  specific column is unknown:
  * `namespace` - (Required in block) Input dataset namespace.
  * `name` - (Required in block) Input dataset name.
  * `field` - (Required in block) Output field this dataset contributes to.
  * `transformation` block — same structure as above.

---

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `process_name` - Full resource name of the Dataplex Process entity created for
  this job. Stable across applies — the same process is reused for all runs.

* `run_name` - Full resource name of the Dataplex Run created for the most recent
  `apply`. Changes on every apply.

* `lineage_event_name` - Full resource name of the Dataplex LineageEvent created
  for the most recent `apply`. Changes on every apply.

* `update_time` - End time of the most recent Run in RFC 3339 format (or start
  time if the run is still in progress).


