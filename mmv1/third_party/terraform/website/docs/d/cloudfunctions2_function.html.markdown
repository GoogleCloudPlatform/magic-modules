---
subcategory: "Cloud Functions (2nd gen)"
page_title: "Google: google_cloudfunctions2_function"
description: |-
  Get information about a Google Cloud Function (2nd gen).
---

# google\_cloudfunctions2\_function

Get information about a Google Cloud Function (2nd gen). For more information see:

* [API documentation](https://cloud.google.com/functions/docs/reference/rest/v2beta/projects.locations.functions).

## Example Usage

```hcl
data "google_cloudfunctions2_function" "my-function" {
  name = "function"
  location = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of a Cloud Function (2nd gen).

* `location` - (Required) The location in which the resource belongs.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `description` - Description of the function.

* `build_config` - Build step of the function that builds a container from the given source.
  Structure is [documented below](#nested_build_config).

* `service_config` - The Service that has been deployed.
  Structure is [documented below](#nested_service_config).

* `event_trigger` - An Eventarc trigger managed by Google Cloud Functions that fires events in response to a condition in another service.
  Structure is [documented below](#nested_event_trigger).

* `labels` - A set of key/value label pairs associated with this Cloud Function.

<a name="nested_build_config"></a>The `build_config` block supports:

* `build` - The Cloud Build name of the latest successful deployment of the function.

* `runtime` - The runtime in which to run the function.

* `entry_point` - The name of the function (as defined in source code) that will be executed.

* `source` - The location of the function source code.
  Structure is [documented below](#nested_source).

* `worker_pool` - Name of the Cloud Build Custom Worker Pool that should be used to build the function.

* `environment_variables` - User-provided build-time environment variables for the function.

* `docker_repository` - User managed repository created in Artifact Registry optionally with a customer managed encryption key.

<a name="nested_source"></a>The `source` block supports:

* `storage_source` - The source location in Google Cloud Storage.
  Structure is [documented below](#nested_storage_source).

* `repo_source` - The source location in a Cloud Source Repository.
  Structure is [documented below](#nested_repo_source).

<a name="nested_storage_source"></a>The `storage_source` block supports:

* `bucket` - Google Cloud Storage bucket containing the source

* `object` - Google Cloud Storage object containing the source.

* `generation` - Google Cloud Storage generation for the object.

<a name="nested_repo_source"></a>The `repo_source` block supports:

* `project_id` - ID of the project that owns the Cloud Source Repository.

* `repo_name` - Name of the Cloud Source Repository.

* `branch_name` - Regex matching branches to build.

* `tag_name` - Regex matching tags to build.

* `commit_sha` - Explicit commit SHA to build.

* `dir` - Directory, relative to the source root, in which to run the build.

* `invert_regex` - Only trigger a build if the revision regex does NOT match the revision regex.

<a name="nested_service_config"></a>The `service_config` block supports:

* `service` - Name of the service associated with a Function.

* `timeout_seconds` - The function execution timeout.

* `available_memory` - The amount of memory available for a function.

* `environment_variables` - Environment variables that shall be available during function execution.

* `max_instance_count` - The limit on the maximum number of function instances that may coexist at a given time.

* `min_instance_count` - The limit on the minimum number of function instances that may coexist at a given time.

* `vpc_connector` - The Serverless VPC Access connector that this cloud function can connect to.

* `vpc_connector_egress_settings` - Available egress settings.

* `ingress_settings` - Available ingress settings.

* `uri` - URI of the Service deployed.

* `gcf_uri` - URIs of the Service deployed

* `service_account_email` - The email of the service account for this function.

* `all_traffic_on_latest_revision` - Whether 100% of traffic is routed to the latest revision.

* `secret_environment_variables` - Secret environment variables configuration.
  Structure is [documented below](#nested_secret_environment_variables).

* `secret_volumes` - Secret volumes configuration.
  Structure is [documented below](#nested_secret_volumes).

<a name="nested_secret_environment_variables"></a>The `secret_environment_variables` block supports:

* `key` - Name of the environment variable.

* `project_id` - Project identifier of the project that contains the secret.

* `secret` - Name of the secret in secret manager.

* `version` - Version of the secret.

<a name="nested_secret_volumes"></a>The `secret_volumes` block supports:

* `mount_path` - The path within the container to mount the secret volume.

* `project_id` - Project identifier of the project that contains the secret.

* `secret` - Name of the secret in secret manager

* `versions` - List of secret versions to mount for this secret.
  Structure is [documented below](#nested_versions).

<a name="nested_versions"></a>The `versions` block supports:

* `version` - Version of the secret.

* `path` - Relative path of the file under the mount path where the secret value for this version will be fetched and made available.

<a name="nested_event_trigger"></a>The `event_trigger` block supports:

* `trigger` - The resource name of the Eventarc trigger.

* `trigger_region` - The region that the trigger will be in.

* `event_type` - The type of event to observe.

* `event_filters` - Criteria used to filter events.
  Structure is [documented below](#nested_event_filters).

* `pubsub_topic` - The name of a Pub/Sub topic in the same project that will be used as the transport topic for the event delivery.

* `service_account_email` - The email of the service account for this function.

* `retry_policy` - Describes the retry policy in case of function's execution failure.

<a name="nested_event_filters"></a>The `event_filters` block supports:

* `attribute` - The name of a CloudEvents attribute.

* `value` - The value for the attribute.

* `operator` - The operator used for matching the events with the value of the filter.
