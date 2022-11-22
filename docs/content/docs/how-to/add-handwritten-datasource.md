---
title: "Add a handwritten datasource"
weight: 22
---

# Add a handwritten datasource

**Note** : only handwritten datasources are currently supported

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
  datasource.
1. Add documentation.

For creating a datasource based off an existing resource you can [make use of the
schema directly](https://github.com/GoogleCloudPlatform/magic-modules/blob/1d293f7bfadacaa20580874c8e8634827fb99a14/mmv1/third_party/terraform/data_sources/data_source_cloud_run_service.go).
Otherwise [implementing the schema directly](https://github.com/GoogleCloudPlatform/magic-modules/blob/1d293f7bfadacaa20580874c8e8634827fb99a14/mmv1/third_party/terraform/data_sources/data_source_google_compute_address.go),
similar to normal resource creation, is the desired path.
