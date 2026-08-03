---
subcategory: "Cloud IAM"
description: |-
  Get the JSON Web Key Set (JWKS) for a Workload Identity Pool.
---

# google_iam_workload_identity_pool_jwks

Get the JSON Web Key Set (JWKS) for a Workload Identity Pool. This data source retrieves public keys used to verify OIDC tokens issued by Google Cloud Workload Identity Pools.

## Example Usage - Basic

```hcl
data "google_iam_workload_identity_pool_openid_config" "oidc" {
  workload_identity_pool_id = "organizations/123456789/locations/global/workloadIdentityPools/agents.global.org-123456789.system.id.goog"
}

data "google_iam_workload_identity_pool_jwks" "example" {
  jwks_uri = data.google_iam_workload_identity_pool_openid_config.oidc.jwks_uri
}

output "jwks_json" {
  value = data.google_iam_workload_identity_pool_jwks.example.jwks_json
}

output "first_key_id" {
  value = data.google_iam_workload_identity_pool_jwks.example.keys[0].kid
}
```

## Argument Reference

The following arguments are supported:

* `jwks_uri` - (Required) The JWKS URI to retrieve the public keys from (e.g. from `google_iam_workload_identity_pool_openid_config.jwks_uri`).

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `jwks_json` - The raw JSON string representation of the JWKS response.

* `keys` - The JWKS for this OP. Structure is documented below.

The `keys` block contains:

* `kty` - Key type. Currently "RSA".

* `use` - Public key use. Currently "sig".

* `kid` - Key ID.

* `n` - Modulus value for kty="RSA".

* `e` - Exponent value for kty="RSA".

* `alg` - Algorithm intended for use with the key. Currently "RS256".
