---
subcategory: "Cloud IAM"
description: |-
  Get JSON Web Key Set (JWKS) public keys for a Workload Identity Pool from GCP.
---

# google_iam_workload_identity_pool_jwks

Get the JSON Web Key Set (JWKS) public keys (`/openid/jwks`) for an Agent Workload Identity Pool from GCP.

## Example Usage

```tf
data "google_iam_workload_identity_pool_openid_config" "oidc" {
  resource_name = "https://sts.googleapis.com/v1/organizations/433637338589/locations/global/workloadIdentityPools/agents.global.org-433637338589.system.id.goog/.well-known/openid-configuration"
}

data "google_iam_workload_identity_pool_jwks" "example" {
  resource_name = data.google_iam_workload_identity_pool_openid_config.oidc.jwks_uri
}
```

## Argument Reference

The following arguments are supported:

* `resource_name` - (Required) The JWKS URI to retrieve the public keys from (e.g. from `google_iam_workload_identity_pool_openid_config.jwks_uri`).

- - -

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `jwks_json` - The raw JSON string representation of the JWKS response.

* `keys` - The list of public keys in the JSON Web Key Set. Structure is [documented below](#nested_keys).

<aByName name="nested_keys"></a>The `keys` block contains:

* `kty` - The key type (e.g. `RSA`).

* `use` - The intended use of the public key (e.g. `sig`).

* `kid` - The unique identifier for the key.

* `n` - The modulus for the RSA public key.

* `e` - The exponent for the RSA public key.

* `alg` - The algorithm intended for use with the key (e.g. `RS256`).
