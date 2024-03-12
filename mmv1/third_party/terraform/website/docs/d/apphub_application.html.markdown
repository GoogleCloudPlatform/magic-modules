---
subcategory: "Apphub"
description: |-
  Application is a functional grouping of Services and Workloads that helps achieve a desired end-to-end business functionality.
---

# google\_apphub\_application

Application is a functional grouping of Services and Workloads that helps achieve a desired end-to-end business functionality. Services and Workloads are owned by the Application.


## Example Usage


```hcl
data "google_apphub_application" "application" {
  project = "project-id"
  application_id = "application"
  location = "location"
}
```

## Argument Reference

The following arguments are supported:

* `project` - The host project of the application.
* `application_id` - (Required) The user-defined identifier of the application.
* `location` - (Required) The location of the application.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `scope` -
  (Required)
  Scope of an application.
  Structure is [documented below](#nested_scope).


<a name="nested_scope"></a>The `scope` block supports:

* `type` -
  (Required)
  Required. Scope Type. 
   Possible values:
  REGIONAL
  Possible values are: `REGIONAL`.

- - -


* `display_name` -
  (Optional)
  Optional. User-defined name for the Application.

* `description` -
  (Optional)
  Optional. User-defined description of an Application.

* `attributes` -
  (Optional)
  Consumer provided attributes.
  Structure is [documented below](#nested_attributes).

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/applications/{{application_id}}`

* `name` -
  Identifier. The resource name of an Application. Format:
  "projects/{host-project-id}/locations/{location}/applications/{application-id}"

* `create_time` -
  Output only. Create time.

* `update_time` -
  Output only. Update time.

* `uid` -
  Output only. A universally unique identifier (in UUID4 format) for the `Application`.

* `state` -
  Output only. Application state. 
   Possible values:
   STATE_UNSPECIFIED
   CREATING
   ACTIVE
   DELETING


<a name="nested_attributes"></a>The `attributes` block supports:

* `criticality` -
  (Optional)
  Criticality of the Application, Service, or Workload
  Structure is [documented below](#nested_criticality).

* `environment` -
  (Optional)
  Environment of the Application, Service, or Workload
  Structure is [documented below](#nested_environment).

* `developer_owners` -
  (Optional)
  Optional. Developer team that owns development and coding.
  Structure is [documented below](#nested_developer_owners).

* `operator_owners` -
  (Optional)
  Optional. Operator team that ensures runtime and operations.
  Structure is [documented below](#nested_operator_owners).

* `business_owners` -
  (Optional)
  Optional. Business team that ensures user needs are met and value is delivered
  Structure is [documented below](#nested_business_owners).


<a name="nested_criticality"></a>The `criticality` block supports:

* `type` -
  (Required)
  Criticality type.
  Possible values are: `MISSION_CRITICAL`, `HIGH`, `MEDIUM`, `LOW`.

<a name="nested_environment"></a>The `environment` block supports:

* `type` -
  (Required)
  Environment type.
  Possible values are: `PRODUCTION`, `STAGING`, `TEST`, `DEVELOPMENT`.

<a name="nested_developer_owners"></a>The `developer_owners` block supports:

* `display_name` -
  (Optional)
  Optional. Contact's name.

* `email` -
  (Required)
  Required. Email address of the contacts.

<a name="nested_operator_owners"></a>The `operator_owners` block supports:

* `display_name` -
  (Optional)
  Optional. Contact's name.

* `email` -
  (Required)
  Required. Email address of the contacts.

<a name="nested_business_owners"></a>The `business_owners` block supports:

* `display_name` -
  (Optional)
  Optional. Contact's name.

* `email` -
  (Required)
  Required. Email address of the contacts.

