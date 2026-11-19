---
subcategory: "Vertex AI"
description: |-
  Get information about a Vertex AI publisher model configuration.
---

# google_vertex_ai_publisher_model_config

Use this data source to get the per-model configuration for a Vertex AI Model Garden publisher model in a given location (regional or `global`).

To get more information about PublisherModelConfig, see:

* [API documentation](https://cloud.google.com/vertex-ai/docs/reference/rest/v1beta1/projects.locations.publishers.models/fetchPublisherModelConfig)
* How-to Guides
    * [Inference metadata logging](https://docs.cloud.google.com/gemini-enterprise-agent-platform/models/capabilities/inference-metadata-logging)
    * [Use models in Model Garden](https://cloud.google.com/vertex-ai/generative-ai/docs/model-garden/use-models)

~> **Warning:** This datasource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](../guides/provider_versions.html.markdown) for more details on beta datasources.

## Example Usage

```hcl
data "google_vertex_ai_publisher_model_config" "config" {
  provider  = google-beta
  region    = "us-central1"
  publisher = "google"
  model     = "gemini-2.5-flash"
}
```

## Argument Reference

The following arguments are supported:

* `publisher` - (Required) The Model Garden publisher identifier, for example `google` or `anthropic`.

* `model` - (Required) The publisher model identifier, for example `gemini-2.5-flash` or `claude-sonnet-4-5`.

- - -

* `region` - (Optional) The location of the publisher model, for example `global` or `us-central1`. If not specified, the provider region is used.

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_vertex_ai_publisher_model_config](https://registry.terraform.io/providers/hashicorp/google-beta/latest/docs/resources/vertex_ai_publisher_model_config#argument-reference) resource for details of all the available attributes.
