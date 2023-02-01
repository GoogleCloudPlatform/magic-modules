subcategory: "Apigee"
page_title: "Google: google_apigee_shared_flow"
description: |-
  You can combine policies and resources into a shared flow that you can consume from multiple API proxies, and even from other shared flows.
---

# google\_apigee\_shared\_flow

You can combine policies and resources into a shared flow that you can consume from multiple API proxies, and even from other shared flows. Although it's like a proxy, a shared flow has no endpoint. It can be used only from an API proxy or shared flow that's in the same organization as the shared flow itself.


To get more information about SharedFlow, see:

* [API documentation](https://cloud.google.com/apigee/docs/reference/apis/apigee/rest/v1/organizations.sharedflows)
* How-to Guides
    * [Sharedflows](https://cloud.google.com/apigee/docs/resources)


## Argument Reference

The following arguments are supported:


* `name` -
  (Required)
  The ID of the shared flow.

* `org_id` -
  (Required)
  The Apigee Organization associated with the Apigee instance,
  in the format `organizations/{{org_name}}`.

* `config_bundle` -
  (Required)
  The configuration bundle zip file path.

- - -



## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `organizations/{{org_id}}/sharedflows/{{name}}`

* `meta_data` -
  Metadata describing the shared flow.
  Structure is [documented below](#nested_meta_data).

* `revision` -
  A list of revisions of this shared flow.

* `latest_revision_id` -
  The id of the most recently created revision for this shared flow.


<a name="nested_meta_data"></a>The `meta_data` block contains:

* `created_at` -
  (Optional)
  Time at which the API proxy was created, in milliseconds since epoch.

* `last_modified_at` -
  (Optional)
  Time at which the API proxy was most recently modified, in milliseconds since epoch.

* `sub_type` -
  (Optional)
  The type of entity described

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import


SharedFlow can be imported using any of these accepted formats:

```
$ terraform import google_apigee_shared_flow.default {{org_id}}/sharedflows/{{name}}
```
