---
title: "Add a datasource"
summary: "Datasources are like terraform resources except they don't *create* anything."
weight: 60
aliases:
  - /docs/how-to/add-handwritten-datasource
  - /how-to/add-handwritten-datasource
  - /docs/how-to/add-handwritten-datasource-documentation
  - /how-to/add-handwritten-datasource-documentation
---

# Add a datasource

**Note:** only handwritten datasources are currently supported.

Datasources are like resources except they don't *create* anything.
They are read-only operations that will expose values needed for subsequent
resource operations.

Most datasources correspond to a resource, and will automatically add fields when they're
added to that resource. You can create a new datasource of this type as follows:

1. Create a new datasource file in the appropriate service folder inside [`magic-modules/mmv1/third_party/terraform/services`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/services). The name of the file should be `data_source_PRODUCT_RESOURCE.go`. Here's an example:
   ```go
   package memcache

   import (
    "fmt"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-provider-google/google/tpgresource"
    transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
   )

   func DataSourceMemcacheInstance() *schema.Resource {
    // Generate datasource schema from resource
    dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceMemcacheInstance().Schema)

    // Set 'Required' schema elements
    tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")
    // Set 'Optional' schema elements
    tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")

    return &schema.Resource{
      Read:   dataSourceMemcacheInstanceRead,
      Schema: dsSchema,
    }
   }

   func dataSourceMemcacheInstanceRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*transport_tpg.Config)

    id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{region}}/instances/{{name}}")
    if err != nil {
      return fmt.Errorf("Error constructing id: %s", err)
    }
    d.SetId(id)

    err = resourceMemcacheInstanceRead(d, meta)
    if err != nil {
      return err
    }

    if err := tpgresource.SetDataSourceLabels(d); err != nil {
      return err
    }

    if d.Id() == "" {
      return fmt.Errorf("%s not found", id)
    }
    return nil
   }
   ```

   Important things to note:

   - `tpgresource.DatasourceSchemaFromResourceSchema` ensures that the datasource schema stays in sync with the resource schema.
   - `tpgresource.AddRequiredFieldsToSchema` and `tpgresource.AddOptionalFieldsToSchema` allow "overriding" whether a specific field is optional or required.
   - The Read function for the datasource is a thin wrapper around the Read function for the related resource. This ensures that new fields on the resource are automatically read for the datasource as well.

1. Create a new test file in the same folder. The name of the file should be `data_source_PRODUCT_RESOURCE_test.go`. Here's an example:

   ```go
   package memcache_test

   import (
    "testing"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
    "github.com/hashicorp/terraform-provider-google/google/acctest"
   )

   func TestAccMemcacheInstanceDatasourceConfig(t *testing.T) {
    t.Parallel()

    context := map[string]interface{}{
      "random_suffix": acctest.RandString(t, 10),
    }

    acctest.VcrTest(t, resource.TestCase{
      PreCheck:                 func() { acctest.AccTestPreCheck(t) },
      ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
      CheckDestroy:             testAccCheckMemcacheInstanceDestroyProducer(t),
      Steps: []resource.TestStep{
        {
          Config: testAccMemcacheInstanceDatasourceConfig(context),
          Check: acctest.CheckDataSourceStateMatchesResourceState(
            "data.google_memcache_instance.default",
            "google_memcache_instance.instance",
          ),
        },
      },
    })
   }

   func testAccMemcacheInstanceDatasourceConfig(context map[string]interface{}) string {
    return acctest.Nprintf(`
   resource "google_compute_network" "memcache_network" {
     name                     = "test-network"
   }

   resource "google_compute_global_address" "service_range" {
     name                     = "address"
     purpose                  = "VPC_PEERING"
     address_type             = "INTERNAL"
     prefix_length            = 16
     network                  = google_compute_network.memcache_network.id
   }

   resource "google_service_networking_connection" "private_service_connection" {
     network                  = google_compute_network.memcache_network.id
     service                  = "servicenetworking.googleapis.com"
     reserved_peering_ranges  = [google_compute_global_address.service_range.name]
   }

   resource "google_memcache_instance" "instance" {
     name                     = "test-instance"
     authorized_network       = google_service_networking_connection.private_service_connection.network
     region                   = "us-central1"
     node_config {
       cpu_count              = 1
       memory_size_mb         = 1024
     }
     node_count               = 1
   }
   data "google_memcache_instance" "default" {
   name                       = google_memcache_instance.instance.name
   region                     = "us-central1"
   }
   `, context)
   }
   ```

   Important things to note:

   - `acctest.CheckDataSourceStateMatchesResourceState` checks that the fields are read properly from the API and will automatically
     stay updated as new fields are added to the resource.
   - If there is `labels` field with type `KeyValueLabels` in the corresponding resource: After calling the resource read method, call the function `tpgresource.SetDataSourceLabels(d)` to make `labels` and `terraform_labels` have all of the labels on the resource.
   - If there is `annotations` field with type `KeyValueAnnotations` in the corresponding resource: After calling the resource read method, call the function `tpgresource.SetDataSourceAnnotations(d)` to make `annotations` have all of the annotations on the resource.

1. Register the datasource to `handwrittenDatasources` in [`magic-modules/mmv1/third_party/terraform/provider/provider_mmv1_resources.go.tmpl`](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/provider/provider_mmv1_resources.go.tmpl)
1. [Add documentation](#add-documentation)

For creating a datasource based off an existing resource you can [make use of the
schema directly](https://github.com/GoogleCloudPlatform/magic-modules/blob/8a8ffc3384a59340f47efe97f18611b6672da9bd/mmv1/third_party/terraform/services/cloudrun/data_source_cloud_run_service.go).
Otherwise [implementing the schema directly](https://github.com/GoogleCloudPlatform/magic-modules/blob/8a8ffc3384a59340f47efe97f18611b6672da9bd/mmv1/third_party/terraform/services/compute/data_source_google_compute_address.go),
similar to normal resource creation, is the desired path.

## Resourceless Datasources

Datasources not backed by a resource are possible to add as well. They follow
the same general steps as adding a resource-based datasource, except that a
full Read method will need to be defined for them rather than calling a
resource's Read method.

Note that while resource-based datasources can depend on the resource read
method for API calls, resourceless datasources need to make them themselves.
An HTTP-based client that's properly configured with logging and retries **must**
be used, such as a client from the https://github.com/googleapis/google-api-go-client
library, or the raw HTTP client used in MMV1 through `SendRequest`.

## Add a field to an existing datasource

If you need to add a field to an existing datasource (that doesn't use `DatasourceSchemaFromResourceSchema`), the process is essentially
the same as for a [handwritten resource]({{< ref "/develop/add-fields" >}}).

## Add documentation

1. Open the data source documentation in [`magic-modules/third_party/terraform/website/docs/d/`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/website/docs/d) using an editor of your choice.
   - The name of the file is the name of the data source without a `google_` prefix. For example, for `google_compute_instance`, the file is called `compute_instance.html.markdown`
2. Modify the documentation as needed according to [Handwritten documentation style guide]({{< ref "/document/handwritten-docs-style-guide" >}}).
   - For resource-based datasources, the "Attribute reference" section should link to the resource documentation so that it doesn't need to be updated as new fields are added. For example:
     ```markdown
     ## Attributes Reference

     See [google_memcache_instance](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/memcache_instance) resource for details of all the available attributes.
     ```
4. [Generate the providers]({{< ref "/develop/generate-providers" >}})
5. Copy and paste the generated documentation into the Hashicorp Registry's [Doc Preview Tool](https://registry.terraform.io/tools/doc-preview) to see how it is rendered.
