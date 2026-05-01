---
subcategory: "Chrome Policy"
page_title: "Google: google_chrome_policies"
description: |-
  Authoritatively manages Chrome policies matching a schema filter for an organizational unit or group.
---

# google_chrome_policies

Provides authoritative management of Chrome policies matching a given schema filter for a specific organizational unit or group. This resource manages the complete set of policies matching the filter — any policies matching the filter that are not defined in the configuration will be removed (reset to inherited).

To get more information, see:

* [API documentation](https://developers.google.com/chrome/policy/reference/rest/v1/customers.policies/resolve)
* How-to Guides
  * [Chrome Policy API Overview](https://developers.google.com/chrome/policy)

~> **Note:** This resource requires the `https://www.googleapis.com/auth/chrome.management.policy` OAuth scope to be added to the provider `scopes` configuration.

~> **Warning:** This resource is authoritative for all policies matching the `schema_filter`. Any policies matching the filter that are not defined in `policies` will be removed (reset to inherited).

## Example Usage

When a single policy is defined, `schema_filter` can be omitted and will be inferred from its `schema`. Because the filter matches only this exact schema, no other policies are affected — this is a safe way to manage an individual policy:

```hcl
resource "google_chrome_policies" "ou_max_connections" {
  org_unit_id = "05qrs456"

  # schema_filter is inferred as "chrome.users.MaxConnectionsPerProxy"

  policies = [
    {
      schema = "chrome.users.MaxConnectionsPerProxy"
      value = {
        maxConnectionsPerProxy = 32
      }
    },
  ]
}
```

When `schema_filter` uses a wildcard, the resource becomes authoritative over all matching policies. The following example force-installs the [Endpoint Verification](https://chromewebstore.google.com/detail/endpoint-verification/callobklhcbilhphinckomhgkigmfocg) extension and grants it access to certificate keys. Because the filter is `chrome.users.apps.*`, this is authoritative over **all** app policies for this org unit — any app policies not declared here will be removed (reset to inherited).

```hcl
resource "google_chrome_policies" "ou_apps" {
  org_unit_id          = "05qrs456"
  schema_filter = "chrome.users.apps.*"

  policies = [
    {
      schema = "chrome.users.apps.InstallType"
      value = {
        appInstallType = "FORCED"
      }
      additional_target_keys = {
        app_id = "chrome:callobklhcbilhphinckomhgkigmfocg"
      }
    },
    {
      schema = "chrome.users.apps.CertificateManagement"
      value = {
        allowAccessToKeys = true
      }
      additional_target_keys = {
        app_id = "chrome:callobklhcbilhphinckomhgkigmfocg"
      }
    },
  ]
}
```

Omitting `policies` or setting it to an empty list ensures that all policies matching the filter are inherited from the parent — useful for enforcing a clean slate on an org unit:

```hcl
resource "google_chrome_policies" "ou_apps_inherited" {
  org_unit_id          = "05qrs456"
  schema_filter = "chrome.users.apps.*"
  policies             = []
}
```

## Argument Reference

The following arguments are supported:

* `customer_id` - (Optional, [ForceNew](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#forcenew)) The ID of the Google Workspace or Cloud Identity customer. Defaults to `my_customer`.

* `org_unit_id` - (Optional, [ForceNew](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#forcenew)) The ID of the organizational unit to apply the policies to. Exactly one of `org_unit_id` or `group_id` must be specified.

* `group_id` - (Optional, [ForceNew](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#forcenew)) The ID of the group to apply the policies to. Exactly one of `org_unit_id` or `group_id` must be specified.

* `schema_filter` - (Optional, Computed, [ForceNew](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#forcenew)) The schema filter defining the authoritative scope of this resource. Can be an exact schema name like `chrome.users.MaxConnectionsPerProxy`, or use a wildcard to match all schemas in a namespace. Wildcards are only supported in the leaf portion of the schema name (e.g. `chrome.users.*`, `chrome.users.apps.*`, `chrome.printers.*`). See the [schema namespaces documentation](https://developers.google.com/chrome/policy/guides/policy-schemas) for details. All policies must have a `schema` that matches the filter. To manage policies in multiple schema namespaces for the same target, use separate `google_chrome_policies` resources. Required when zero or multiple policies are defined. If exactly one policy is defined, it is inferred from that policy's `schema`.

* `policies` - (Optional) A list of policies to enforce. The list is treated as a set — the order of entries is not significant. Each entry is an object with the following fields:

  * `schema` - (Required) The fully qualified name of the policy schema, e.g. `chrome.users.apps.InstallType`.

  * `value` - (Required) The policy values as a map. Accepts native HCL types (strings, booleans, numbers, lists).

  * `additional_target_keys` - (Optional) A map of additional target keys required by some policies to identify a specific target. For example, app-scoped policies under `chrome.users.apps.*` require an `app_id` to specify which extension or app the policy applies to: `{ app_id = "chrome:callobklhcbilhphinckomhgkigmfocg" }`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource in the format `{customer_id}/{orgunits|groups}/{target_id}/{schema_filter}`.

## Import

Chrome Policies can be imported using the format `{customer_id}/{orgunits|groups}/{target_id}/{schema_filter}`:

```
$ terraform import google_chrome_policies.default my_customer/orgunits/03ph8a2z2i02hdg/chrome.users.apps.*
$ terraform import google_chrome_policies.default C0abc123/groups/01abc456/chrome.devices.*
```

Use `my_customer` as the customer ID to refer to the customer associated with the authenticated account.
