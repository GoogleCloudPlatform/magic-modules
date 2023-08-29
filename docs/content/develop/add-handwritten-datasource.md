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

**Note:** only handwritten datasources are currently supported

Datasources are like terraform resources except they don't *create* anything.
They are simply read-only operations that will expose some sort of values needed
for subsequent resource operations. If you're adding a field to an existing
datasource, check the [Resource](#resource) section. Everything there will
be mostly consistent with the type of change you'll need to make. For adding
a new datasource there are 5 steps to doing so.

1. Create a new datasource declaration file and a corresponding test file
1. Add Schema and Read operation implementation
1. Add the datasource to the `provider.go.erb` index
1. Implement a test which will create and resources and read the corresponding
  datasource
1. [Add documentation](#add-documentation)

For creating a datasource based off an existing resource you can [make use of the
schema directly](https://github.com/GoogleCloudPlatform/magic-modules/blob/1d293f7bfadacaa20580874c8e8634827fb99a14/mmv1/third_party/terraform/data_sources/data_source_cloud_run_service.go).
Otherwise [implementing the schema directly](https://github.com/GoogleCloudPlatform/magic-modules/blob/1d293f7bfadacaa20580874c8e8634827fb99a14/mmv1/third_party/terraform/data_sources/data_source_google_compute_address.go),
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

## Add documentation

1. Open the data source documentation in [`magic-modules/third_party/terraform/website/docs/d/`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/website/docs/d) using an editor of your choice.
   - The name of the file is the name of the data source without a `google_` prefix. For example, for `google_compute_instance`, the file is called `compute_instance.html.markdown`
2. Modify the documentation as needed according to [Handwritten documentation style guide]({{< ref "/develop/handwritten-docs-style-guide" >}}).
4. [Generate the providers]({{< ref "/get-started/generate-providers.md" >}})
5. Copy and paste the generated documentation into the Hashicorp Registry's [Doc Preview Tool](https://registry.terraform.io/tools/doc-preview) to see how it is rendered.