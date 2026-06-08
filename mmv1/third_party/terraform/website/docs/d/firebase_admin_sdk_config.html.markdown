---
subcategory: "Firebase"
description: |-
  A Google Cloud Firebase Admin SDK configuration
---

# google_firebase_admin_sdk_config

A Google Cloud Firebase Admin SDK configuration

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](../guides/provider_versions.html.markdown) for more details on beta resources.

To get more information about AdminSdkConfig, see:

* [API documentation](https://firebase.google.com/docs/reference/firebase-management/rest/v1beta1/projects/getAdminSdkConfig)
* How-to Guides
    * [Official Documentation](https://firebase.google.com/)


## Argument Reference
The following arguments are supported:

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `project` -
  The ID of the project in which the resource belongs.

* `database_url` -
  The default Firebase Realtime Database URL.

* `storage_bucket` -
  The default Cloud Storage for Firebase storage bucket name.

* `location_id` -
  The ID of the project's default GCP resource location.
