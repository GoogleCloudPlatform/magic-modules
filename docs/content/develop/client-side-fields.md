---
title: "Client-side fields"
weight: 400
---

# Client-side fields

Client-side fields are most often used as flags to modify the behavior of a Terraform resource. Because they don't correspond to an API field, there are some additional considerations in terms of how to implement them.

Common client-side fields include:

- [`deletion_protection`]({{< ref "/best-practices#deletion_protection" >}})
- [`deletion_policy`]({{< ref "/best-practices#deletion_policy" >}})

{{< tabs "schema" >}}
{{< tab "MMv1" >}}
## Add to the schema

Instead of adding the field in `parameters` or `properties`, use a section called `virtual_fields`.

Example:
```yaml
virtual_fields:
  - !ruby/object:Api::Type::Boolean
    name: 'deletion_protection'
    default_value: true
    description: |
      Whether Terraform will be prevented from destroying the CertificateAuthority.
      When the field is set to true or unset in Terraform state, a `terraform apply`
      or `terraform destroy` that would delete the CertificateAuthority will fail.
      When the field is set to false, deleting the CertificateAuthority is allowed.
```

This will automatically ensure that the field works as users expect.
{{< /tab >}}
{{< tab "Handwritten" >}}
## Add to the schema

Add the field to the schema as usual.

Example:

```go
"deletion_protection": {
	Type:        schema.TypeBool,
	Optional:    true,
	Default:     true,
	Description: `Whether Terraform will be prevented from destroying the instance. When the field is set to true or unset in Terraform state, a terraform apply or terraform destroy that would delete the table will fail. When the field is set to false, deleting the table is allowed.`,
},
```
## Set on read

For fields with default values, you need to explicitly set client-side fields in the Read function to avoid a diff that "sets" the field to its default value when users upgrade to the version that contains the new field.

Example:

```go
// Explicitly set client-side fields to default values if unset
if _, ok := d.GetOkExists("deletion_protection"); !ok {
	if err := d.Set("deletion_protection", true); err != nil {
		return fmt.Errorf("Error setting deletion_protection: %s", err)
	}
}
```

## Short-circuit updates if only client-side fields were modified

Client-side fields can always be updated in-place. However, if the resource is otherwise immutable, you will need to ensure that an Update function is present for the resource. Terraform automatically updates the state based on the plan; this does not need to happen in the Update function unless the value for a field might have changed (for example, based on an API response, which doesn't apply to client-side fields).

If only client-side fields were modified, you can short-circuit the Update function to avoid sending an API request. This is important because the update request will be empty (which causes errors for some APIs.) This can go at the top of the Update function:

```go
clientSideFields := map[string]bool{"deletion_protection": true}
clientSideOnly := true
for field := range ResourceSpannerInstance().Schema {
	if d.HasChange(field) && !clientSideFields[field] {
		clientSideOnly = false
		break
	}
}
if clientSideOnly {
	return nil
}
```

Replace `ResourceSpannerInstance` with the appropriate resource function.
{{< /tab >}}
{{< /tabs >}}

## Implement logic

At this point, you should be ready to implement your logic!
