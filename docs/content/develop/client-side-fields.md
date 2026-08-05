---
title: "Client-side fields"
weight: 150
---

# Client-side fields

Client-side fields are most often used as flags to modify the behavior of a Terraform resource. Because they don't correspond to an API field, there are some additional considerations in terms of how to implement them.

{{% tabs "schema" %}}
{{< tab "MMv1" >}}
## Add to the schema

Instead of adding the field in `parameters` or `properties`, use a section called `virtual_fields`.

Example:
```yaml
virtual_fields:
  - name: 'send_propagated_connection_limit_if_zero'
    type: Boolean
    default_value: false
    description: |
      Controls the behavior of propagated_connection_limit.
      When false, setting propagated_connection_limit to zero causes the provider to use to the API's default value.
      When true, the provider will set propagated_connection_limit to zero.
      Defaults to false.
```

{{< /tab >}}
{{< tab "Handwritten" >}}
## Add to the schema

Add the field to the schema as usual.

Example:

```go
"send_propagated_connection_limit_if_zero": {
  Type:     schema.TypeBool,
  Optional: true,
  Description: `Controls the behavior of propagated_connection_limit.
When false, setting propagated_connection_limit to zero causes the provider to use to the API's default value.
When true, the provider will set propagated_connection_limit to zero.
Defaults to false.`,
  Default: false,
},
```
## Set on read

For fields with default values, you need to explicitly set client-side fields in the Read function to avoid a diff that "sets" the field to its default value when users upgrade to the version that contains the new field.

Example:

```go
// Explicitly set client-side fields to default values if unset
if _, ok := d.GetOkExists("send_propagated_connection_limit_if_zero"); !ok {
	if err := d.Set("send_propagated_connection_limit_if_zero", false); err != nil {
		return fmt.Errorf("Error setting send_propagated_connection_limit_if_zero: %s", err)
	}
}
```

## Short-circuit updates if only client-side fields were modified

Client-side fields can always be updated in-place. However, if the resource is otherwise immutable, you will need to ensure that an Update function is present for the resource. Terraform automatically updates the state based on the plan; this does not need to happen in the Update function unless the value for a field might have changed (for example, based on an API response, which doesn't apply to client-side fields).

If only client-side fields were modified, you can short-circuit the Update function to avoid sending an API request. This is important because the update request will be empty (which causes errors for some APIs.) This can go at the top of the Update function:

```go
clientSideFields := map[string]bool{"send_propagated_connection_limit_if_zero": true}
clientSideOnly := true
for field := range ResourceComputeServiceAttachment().Schema {
	if d.HasChange(field) && !clientSideFields[field] {
		clientSideOnly = false
		break
	}
}
if clientSideOnly {
	return nil
}
```

Replace `ResourceComputeServiceAttachment` with the appropriate resource function.
{{< /tab >}}
{{% /tabs %}}

## Update data source

If the resource has a corresponding data source that calls the resource's Read function, you will need to make the following changes to the data source:

1. Add the client-side field to the data source's Schema as an output-only field. (This will happen automatically for data sources that use `tpgresource.DatasourceSchemaFromResourceSchema`.)

   ```go
   "send_propagated_connection_limit_if_zero": {
     Type:        schema.TypeBool,
     Computed:    true,
   },
   ```
2. Unset the field in the data source's Read function.

   ```go
   if err := d.Set("send_propagated_connection_limit_if_zero", nil); err != nil {
     return fmt.Errorf("Error setting send_propagated_connection_limit_if_zero: %s", err)
   }
   ```

## Implement logic

At this point, you should be ready to implement your logic! For example, the `send_propagated_connection_limit_if_zero` field used here flags if an explicit 0 should be sent to the API for the `propagated_connection_limit` field.

{{% tabs "implementation" %}}
{{< tab "MMv1" >}}
Add the following as [encoder (and update_encoder) custom code]({{< ref "/develop/custom-code#encoder" >}}).

```go
propagatedConnectionLimitProp := d.Get("propagated_connection_limit")
if sv, ok := d.GetOk("send_propagated_connection_limit_if_zero"); ok && sv.(bool) {
  if v, ok := d.GetOkExists("propagated_connection_limit"); ok || !reflect.DeepEqual(v, propagatedConnectionLimitProp) {
    obj["propagatedConnectionLimit"] = propagatedConnectionLimitProp
  }
}

return obj, nil
```
{{< /tab >}}
{{< tab "Handwritten" >}}
Add the following encoders that are referenced within the Create and Update functions respectively, prior to the resource `obj` being sent as a request to the API.
```go
func resourceComputeServiceAttachmentEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
  propagatedConnectionLimitProp := d.Get("propagated_connection_limit")
  if sv, ok := d.GetOk("send_propagated_connection_limit_if_zero"); ok && sv.(bool) {
    if v, ok := d.GetOkExists("propagated_connection_limit"); ok || !reflect.DeepEqual(v, propagatedConnectionLimitProp) {
      obj["propagatedConnectionLimit"] = propagatedConnectionLimitProp
    }
  }

  return obj, nil
}

func resourceComputeServiceAttachmentUpdateEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
  propagatedConnectionLimitProp := d.Get("propagated_connection_limit")
  if sv, ok := d.GetOk("send_propagated_connection_limit_if_zero"); ok && sv.(bool) {
    if v, ok := d.GetOkExists("propagated_connection_limit"); ok || !reflect.DeepEqual(v, propagatedConnectionLimitProp) {
      obj["propagatedConnectionLimit"] = propagatedConnectionLimitProp
    }
  }

  return obj, nil
}
```
{{< /tab >}}
{{% /tabs %}}
