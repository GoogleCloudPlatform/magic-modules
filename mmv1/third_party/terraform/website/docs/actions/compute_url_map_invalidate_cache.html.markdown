---
subcategory: "Compute Engine"
description: |-
    Invalidates Cloud CDN cached content for a global URL map
---

# google_compute_url_map_invalidate_cache (action)

Invalidates Cloud CDN cached content for a global URL map. This is equivalent to calling
[`urlMaps.invalidateCache`](https://docs.cloud.google.com/compute/docs/reference/rest/v1/urlMaps/invalidateCache)
or running `gcloud compute url-maps invalidate-cdn-cache`. The action waits for the invalidation operation 
to complete.

Invalidation does not change the URL map or its backends, so this action does not cause drift. Invalidating
the same content more than once is harmless, so a failed invocation can be retried.

Actions require Terraform 1.14 or later. For an overview of how actions work in the Terraform 
Provider for Google Cloud, see
[Use actions with Terraform Provider for Google Cloud](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/using_actions).

To get more information about cache invalidation, see:

* [API documentation](https://cloud.google.com/compute/docs/reference/rest/v1/urlMaps/invalidateCache)
* How-to Guides
    * [Cache invalidation overview](https://cloud.google.com/cdn/docs/cache-invalidation-overview)
    * [Invalidate cached content](https://cloud.google.com/cdn/docs/invalidating-cached-content)

## Example Usage - Invoke from the CLI

```hcl
resource "google_storage_bucket" "static" {
  name                        = "my-static-site-bucket"
  location                    = "US"
  uniform_bucket_level_access = true
}

resource "google_compute_backend_bucket" "static" {
  name         = "static-backend"
  bucket_name  = google_storage_bucket.static.name
  enable_cdn   = true
}

resource "google_compute_url_map" "static" {
    name            = "static-url-map"
    default_service = google_compute_backend_bucket.static.id
}

action "google_compute_url_map_invalidate_cache" "example" {
  config {
    url_map = google_compute_url_map.static.id
    path    = "/*"
  }
}
```

Run the following command to invalidate all cached content:

```shell
terraform apply -invoke=action.google_compute_url_map_invalidate_cache.example
```

## Example Usage - Invalidate when content changes

This example invalidates the cached copy of `index.html` whenever the object changes. It uses the URL map
from the previous example.

```hcl
resource "google_storage_bucket_object" "index" {
  name   = "index.html"
  bucket = google_storage_bucket.static.name
  source = "${path.module}/site/index.html"
}

action "google_compute_url_map_invalidate_cache" "index" {
  config {
    url_map = google_compute_url_map.static.id
    path    = "/index.html"
  }
}

resource "terraform_data" "index_version" {
    input = google_storage_bucket_object.index.md5hash

    lifecycle {
        action_trigger {
            events  = [after_update]
            actions = [action.google_compute_url_map_invalidate_cache.index]
        }
    }
}
```

## Example Usage - Invalidate by host and cache tags

```hcl
action "google_compute_url_map_invalidate_cache" "product_images" {
  config {
    url_map    = "projects/my-project/global/urlMaps/static-url-map"
    path       = "/images/*"
    host       = "www.example.com"
    cache_tags = ["product-123", "product-456"]
    timeout    = "5m"
  }
}
```

## Argument Reference

The following arguments are supported in the `config` block:

* `url_map` - (Required) The global URL map whose cache to invalidate. Accepts any of the following formats:
    * The URL map name, for example `static-url-map`.
    * A partial URI: `global/urlMaps/{name}` or `projects/{project}/global/urlMaps/{name}`.
    * A self link, such as the value of `google_compute_url_map.static.self_link`.
  
  Regional URL maps are not supported.
  
* `path` - (Required) The path to invalidate. The path must start with `/` and must not contain `?` or `#`,
 because query strings and fragments are not matched. A `*` is allowed only as the final character and
 matches every path with that prefix: `/images/*` invalidates everything under `/images/`, and `/*`
 invalidates all content.

* `host` - (Optional) If set, only content cached for requests with this `Host` header is invalidated. The
 leftmost label can be a wildcard, for example `*.example.com`. If it is not set, the path is invalidated
 for all hosts.

* `cache_tags` - (Optional) Up to 10 cache tags, each up to 120 characters. Only cached objects that have at
 least one of these tags, and that also match `path` (and `host`, if set), are invalidated. Tags are set by 
 the origin with the `Cache-Tag` response header. To invalidate by tag alone, set `path` to `/*`.

* `project` - (Optional) The ID of the project that contains the URL map. If it is not provided,
 the project in `url_map` is used, and then the provider project. If you set both `project` and a 
 project in `url_map`, they must match.

* `timeout` - (Optional) The maximum time to wait for the invalidation operation to complete, as a duration
 string such as `"5m"`. The value must be between `"1s"` and `"1h"`. Defaults to `"20m"`.

## Behavior

* The action submits the invalidation and waits for the operation to complete. An invalidation usually
 takes effect in about 10 seconds. In rare cases, Cloud CDN reports an invalidation as complete before
 every cache has processed it.
* If the URL map does not exist, the action returns an error.
* Each invocation submits a new invalidation. Cloud CDN limits how many invalidation requests a project can
 submit; see [Cache invalidation overview](https://cloud.google.com/cdn/docs/cache-invalidation-overview).
 Attach `action_trigger` to a resource that changes only when content changes, rather than to one that is
 updated on every apply.
