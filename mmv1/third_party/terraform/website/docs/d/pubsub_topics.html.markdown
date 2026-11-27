---
subcategory: "Cloud Pub/Sub"
description: |-
  Get a list of the Pub/Sub topics in a project.
---

# google_pubsub_topics

Get a list of the Pub/Sub topics in a project. For more information see
[the official documentation](https://cloud.google.com/pubsub/docs/list-topics)
and [API](https://cloud.google.com/pubsub/docs/reference/rest/v1/projects.topics/list).

## Example Usage

```hcl
data "google_pubsub_topics" "all" {
  project = "my-project"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The ID of the project whose topics are listed. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/topics`

* `topics` - The list of topics in the project.

See [google_pubsub_topic](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/pubsub_topic) resource for details of the available attributes on each topic, with these differences:

* `id` is the identifier of the topic, with format `projects/{{project}}/topics/{{name}}`.
* `labels` holds every label of the topic.
* The `project`, `deletion_policy`, `tags`, `terraform_labels` and `effective_labels` arguments and attributes of that resource are not part of each topic.
