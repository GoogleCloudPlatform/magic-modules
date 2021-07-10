---
subcategory: "Firebase"
layout: "google"
page_title: "Google: google_firebase_ios_app"
sidebar_current: "docs-google-firebase-ios-app"
description: |-
  A Google Cloud Firebase App for iOS
---

# google\_firebase\_ios\_app

Provides access to an iOS app's attributes within Google Cloud Firebase. For more information, see the [official documentation](https://firebase.google.com/docs/ios/setup) and [API](https://firebase.google.com/docs/projects/api/reference/rest/v1beta1/projects.iosApps)

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.

## Argument Reference

The following arguments are supported:

* `app_id` - (Required) Immutable. The globally unique, Firebase-assigned identifier of the App. This identifier should be treated as an opaque token, as the data format is not specified.

- - -

* `project` - (Optional) The parent Firebase project's ProjectNumber (recommended) or its ProjectId.

## Attributes Reference

The following attributes are exported:

* `name` - The fully qualified resource name of the App, for example: `projects/{{project}}/iosApps/{{app_id}}`
* `app_id` - Immutable. The globally unique, Firebase-assigned identifier of the App. This identifier should be treated as an opaque token, as the data format is not specified.
* `display_name` - The user-assigned display name for the IosApp.
* `project_id` - Immutable. A user-assigned unique identifier of the parent FirebaseProject for the IosApp.
* `bundle_id` - Immutable. The canonical bundle ID of the iOS app as it would appear in the iOS AppStore.
* `app_store_id` - The automatically generated Apple ID assigned to the iOS app by Apple in the iOS App Store.
