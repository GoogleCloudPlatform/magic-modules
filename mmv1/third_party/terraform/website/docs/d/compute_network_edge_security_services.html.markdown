---
subcategory: "Compute Engine"
page_title: "Google: google_compute_network_edge_security_services"
description: |-
  Represents a Google Cloud Armor network edge security service resource.
---

# google\_compute\_network\_edge\_security\_services

Represents a Google Cloud Armor network edge security service resource.

see the [official documentation](https://cloud.google.com/armor/docs/configure-security-policies)
and the [API](https://cloud.google.com/compute/docs/reference/rest/v1/networkEdgeSecurityServices).

## Example Usage

```hcl
resource "google_compute_network_edge_security_services" "services" {
  name = "my-policy"
  description = "default rule"
  }
}
```

## Example Security Police

```hcl
resource "google_compute_network_edge_security_services" "services" {
    name        = "%s"
    description = "basic network edge security services"
    security_policy = google_compute_region_security_policy.policy.self_link
}

resource "google_compute_region_security_policy" "policy" {
    name        = "%s"
    description = "default rule"
    type = "CLOUD_ARMOR_NETWORK"

    ddos_protection_config {
        ddos_protection = "ADVANCED"
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the security policy.

* `description` - (Optional) An optional description of this security policy. Max size is 2048.

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

* `rule` - (Optional) The set of rules that belong to this policy. There must always be a default
    rule (rule with priority 2147483647 and match "\*"). If no rules are provided when creating a
    security policy, a default rule with action "allow" will be added. Structure is [documented below](#nested_rule).

* `region` - [Output Only] URL of the region where the resource resides. You must specify this field as part of the HTTP request URL. It is not settable as a field in the request body.

* `securityPolicy` - The resource URL for the network edge security service associated with this network edge security service. Structure is [documented below](#nested_rule).

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - an identifier for the resource with format `projects/{{project}}/regions/{{region}}/networkEdgeSecurityServices/{{name}}`

* `fingerprint` - Fingerprint of this resource.

* `self_link` - The URI of the created resource.
