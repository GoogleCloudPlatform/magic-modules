---
subcategory: "Apikeys"
description: |-
  Produces an ephemeral resource for an API key string
---

# google_apikeys_key

This ephemeral resource provides access to the key string for an existing API Keys API key. It can be used to pass an API key into write-only arguments without storing the key string in Terraform state.

To get more information about API Keys, see:

* [API documentation](https://cloud.google.com/api-keys/docs/reference/rest/v2/projects.locations.keys/getKeyString)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/api-keys/docs)

## Example Usage

```hcl
resource "google_apikeys_key" "primary" {
  name         = "my-key"
  display_name = "sample-key"
}

ephemeral "google_apikeys_key" "primary" {
  name = google_apikeys_key.primary.id
}

resource "google_secret_manager_secret_version" "api_key" {
  secret                 = google_secret_manager_secret.api_key.id
  secret_data_wo         = ephemeral.google_apikeys_key.primary.key_string
  secret_data_wo_version = "1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The API key resource name or key id. This can be a full resource name in the format `projects/{{project}}/locations/global/keys/{{key}}`, or the final key id when `project` is set or available from provider configuration.

- - -

* `project` - (Optional) The project to get the API key for. If it is not provided, the provider project is used. This field is inferred when `name` is a full resource name.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `key_string` - (Output) The encrypted and signed value held by this API key.
* `id` - (Output) The full API key resource name. Format: `projects/{{project}}/locations/global/keys/{{key}}`.
* `project` - (Output) The project the API key belongs to.
