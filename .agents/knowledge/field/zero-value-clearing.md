---
name: zero-value-clearing
description: When updateMask suffices vs. when send_empty_value, allow_empty_object, or GetRawConfigAt is needed.
topics: [field, custom-code]
task_types: [field-add, new-resource, bug-fix]
source: docs/content/reference/field.md + mmv1/templates/terraform/resource.go.tmpl
status: approved
last_verified: 2026-11-06
---

# Zero-Value Clearing (`false`, `0`, `""`, `{}`) & Raw Config Access

## 1. When `send_empty_value: true` Is (and Is Not) Needed

Standard Google AIP-134 `PATCH` endpoints clear implicit-presence `proto3` scalar fields when MMv1 adds the changed field to `?updateMask=...` and omits the zero value from the JSON body. **Do not** add `send_empty_value: true` to ordinary `Optional` primitives by default.

Set `send_empty_value: true` (see [`docs/content/reference/field.md#send_empty_value`](../../../docs/content/reference/field.md#send_empty_value)) **only** when:
1. **Non-zero server default on `Create`:** Omitting `false` or `0` on `POST` causes the server to apply a non-zero default (e.g., `true`).
2. **Explicit-presence `REQUIRED` scalar:** `proto3 optional` + `(google.api.field_behavior) = REQUIRED` rejects omitted zero values (`false`/`0`).
3. **Full-replacement `PUT` or explicit-presence `PATCH`:** Verified API behavior or wire logs show mask-only omission fails to reset a prior value.
4. **Empty presence object (`field_name {}`):** Paired with `allow_empty_object: true` on a `NestedObject` when the API requires `{}` on the wire (watch for [tpg#13201](https://github.com/hashicorp/terraform-provider-google/issues/13201), which sends `"field": null` when a top-level `NestedObject` with `send_empty_value: true` is unset).

When `send_empty_value: true` or custom zero-value clearing is required for `PATCH` resets, include a `non-zero -> zero` transition in the update test.

## 2. Nested Fields with `default_from_api` (`rd.GetRawConfigAt`)

Top-level `send_empty_value: true` + `default_from_api: true` omits unset fields automatically, but nested fields inside objects require a `custom_expand` returning a typed pointer (`*bool` / `*int64`, serialized as `false`/`0`) when configured in HCL and `nil` when unset so the API default applies:

```go
func expand{{$.GetPrefix}}{{$.TitlelizeProperty}}(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	b, ok := v.(bool)
	if !ok {
		return nil, nil
	}
	if b {
		return b, nil
	}
{{- if not (contains $.ResourceMetadata.ProductMetadata.Compiler "terraformgoogleconversion") }}
	if rd, ok := d.(*schema.ResourceData); ok {
		path := cty.GetAttrPath("parent_block").IndexInt(0).GetAttr("{{ underscore $.Name }}")
		if val, _ := rd.GetRawConfigAt(path); !val.IsNull() && val.IsKnown() {
			f := false
			return &f, nil
		}
	}
{{- end }}
	return nil, nil
}
```

Live examples: `mmv1/templates/terraform/custom_expand/compute_backend_service_cdn_policy_ttl.go.tmpl`, `mmv1/templates/terraform/custom_expand/modelarmor_template_data_residency_compliant.go.tmpl`. Inside `CustomizeDiff` (`*schema.ResourceDiff`), inspect `diff.GetRawConfig()`.

## Do NOT Use `send_empty_value: true` For

- Standard `Optional` primitives whose `PATCH` + `updateMask` already clears omitted zero values, `immutable: true` fields with zero defaults, optional strings where `""` triggers `400` validation errors, or `output: true` fields.
