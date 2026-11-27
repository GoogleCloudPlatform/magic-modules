---
subcategory: "Cloud Pub/Sub"
description: |-
  Get a list of the Pub/Sub subscriptions in a project.
---

# google_pubsub_subscriptions

Get a list of the Pub/Sub subscriptions in a project. For more information see
[the official documentation](https://cloud.google.com/pubsub/docs/list-subscriptions)
and [API](https://cloud.google.com/pubsub/docs/reference/rest/v1/projects.subscriptions/list).

## Example Usage

```hcl
data "google_pubsub_subscriptions" "all" {
  project = "my-project"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The ID of the project whose subscriptions are listed. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/subscriptions`

* `subscriptions` - The list of subscriptions in the project.

See [google_pubsub_subscription](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/pubsub_subscription) resource for details of the available attributes on each subscription, with these differences:

* `id` is the identifier of the subscription, with format `projects/{{project}}/subscriptions/{{name}}`.
* `topic` is the full name of the topic, with format `projects/{{project}}/topics/{{name}}`, or `_deleted-topic_` if the topic has been deleted.
* `labels` holds every label of the subscription.
* The `project`, `deletion_policy`, `tags`, `terraform_labels` and `effective_labels` arguments and attributes of that resource are not part of each subscription.
