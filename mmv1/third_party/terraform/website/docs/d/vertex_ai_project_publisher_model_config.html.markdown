---
subcategory: "Vertex AI"
description: |-
  Get information about the project-wide Vertex AI publisher model configuration for a Google Cloud project.
---

# google_vertex_ai_project_publisher_model_config

Use this data source to get the project-wide default configuration for Vertex AI Model Garden publisher models in a Google Cloud project.

To get more information about ProjectPublisherModelConfig, see:

* [API documentation](https://docs.cloud.google.com/gemini-enterprise-agent-platform/reference/rest/v1beta1/projects/fetchPublisherModelConfig)
* How-to Guides
    * [Inference metadata logging](https://docs.cloud.google.com/gemini-enterprise-agent-platform/models/capabilities/inference-metadata-logging)
    * [Use models in Model Garden](https://cloud.google.com/vertex-ai/generative-ai/docs/model-garden/use-models)

~> **Warning:** This datasource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](../guides/provider_versions.html.markdown) for more details on beta datasources.

## Example Usage

```hcl
data "google_vertex_ai_project_publisher_model_config" "config" {
  provider = google-beta
  project  = "my-project-id"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_vertex_ai_project_publisher_model_config](https://registry.terraform.io/providers/hashicorp/google-beta/latest/docs/resources/vertex_ai_project_publisher_model_config#argument-reference) resource for details of all the available attributes.
