---
title: "Deletion policy"
weight: 20
aliases:
- /best-practices/deletion-behaviors
---
# Deletion policy

This page documents how to work with the `deletion_policy` field present in nearly all provider resources.

{{% tabs "schema" %}}
{{< tab "MMv1" >}}
MMv1 resources support `deletion_policy` out of the box without any complication or changes needed by a contributor.

There are three values allowed on every `deletion_policy` field:
- `PREVENT`: This will block deletion of the resource.
- `ABANDON`: Deletion of the resource will remove it from Terraform state but not from the API.
- `DELETE`: Deletion of the resource will remove it from Terraform state and from the API.

## Updating the default value of `deletion_policy` for a resource
All resources will use a value of "DELETE" by default if unspecified, but this default can be changed to another value by including the following:
```yaml
deletion_policy_default: "PREVENT" #this can be any string
```
We recommend setting the default value to `PREVENT` for resources where it's important to protect against accidental deletion.

## Supporting a custom value for deletion_policy
If implementing support for additional values such as "FORCE", the following steps can be taken.
Add a `pre_delete` constant to the `custom_code` block in the resource YAML file, that performs the logic with the corresponding value. Example implemention here: [Example #1](https://github.com/GoogleCloudPlatform/magic-modules/blob/03e57f68a1f8fa32923d87c224049fb3ac4802e1/mmv1/products/datastream/PrivateConnection.yaml#L49), [Example #2](https://github.com/GoogleCloudPlatform/magic-modules/blob/c5ef760d4e089ca073ca1e40134dda60a7b81096/mmv1/templates/terraform/pre_delete/private_connection.go.tmpl#L1)

Add [acceptance tests](https://googlecloudplatform.github.io/magic-modules/test/test/#add-an-acceptance-tests) for any custom values to verify their implementation works and does not disrupt the usage of the global values "PREVENT", "ABANDON", and "DELETE".

In the resource YAML, set `deletion_policy_custom_docs: true` to prevent the default universal deletion policy field documentation from being generated. Then, add a virtual field with the `exclude: true` attribute set to generate documentation. For example:
```yaml
deletion_policy_custom_docs: true
virtual_fields:
  - name: 'deletion_policy'
    description: |
      The deletion policy for the private connection. Setting `FORCE` will also delete any child
      routes that belong to this private connection. Setting `DEFAULT` will fail the delete if
      child routes exist. Defaults to `FORCE` for backwards compatibility.

      When a 'terraform destroy' or 'terraform apply' would delete the resource,
      the command will fail if this field is set to "PREVENT" in Terraform state.
      When set to "ABANDON", the command will remove the resource from Terraform
      management without updating or deleting the resource in the API.
      When set to "DELETE", the command will behave as if set to "DEFAULT".
    type: String
    exclude: true
```

## Excluding a resource from universal deletion_policy
If a resource is incompatible with `deletion_policy` the following can be added to the resource YAML file, and all related code will not be generated:
```yaml
deletion_policy_exclude: true
```
{{< /tab >}}
{{< tab "Handwritten" >}}
All handwritten resources need to support `deletion_policy` unless deemed incompatible. The following code snippets need to be included in a resource to support this. If following steps under [Add Resource]({{< ref "/develop/add-resource/#handwritten" >}}), they will be included in your resource go file already. 

### CustomizeDiff
Add `tpgresource.DefaultProviderDeletionPolicy()` to the CustomDiff attribute of a given resource's `*schema.Resource`. If no existing CustomizeDiff is present, the whole attribute will need to be added.
```go
CustomizeDiff: customdiff.All(
            tpgresource.DefaultProviderDeletionPolicy("DELETE"),
        ),
```
If the default is being changed, update it from "DELETE" here.

### Schema
Add following to the top level schema of the resource:
```go
"deletion_policy": tpgresource.DeletionPolicySchemaEntry("DELETE"),
```
If the default is being changed, update it from "DELETE" here.

### Read
Add the following to the end of a resource's Read() function:
```go
    if err := tpgresource.DeletionPolicyReadDefault(d, config, "DELETE"); err != nil {
        return err
    }
```
If the default is being changed, update it from "DELETE" here.

### Update
Add the following to the start of a resource's Update() function:
```go
if tpgresource.DeletionPolicyPreUpdate(d, RESOURCENAME) {
    return RESOURCENAME().Read(d, meta)
}
```
If the resource does not support updating, implement the following Update() function for the resource:
```go
//UDP update start
func resourceRESOURCENAMEUpdate(d *schema.ResourceData, meta interface{}) error {
    // Only the root field "deletion_policy", "labels", "terraform_labels", and virtual fields are mutable
    return resourceRESOURCENAMERead(d, meta)
}
//UDP update end
```

### Delete
Add following to the start of a resource's Delete() function, modifying as necessary for what custom values are supported:
```go
if ok, err := tpgresource.DeletionPolicyPreDelete(d); err != nil{
    return err
}else if ok{
    return nil
}
```

### Metadata
Add following line to the resource's meta.yaml:
```yaml
- field: 'deletion_policy'
  provider_only: true
```

### Docs
Add the following to the resource's documentation markdown file, modifying as necessary depending for what custom values are supported:
```markdown
* `deletion_policy` - (Optional) Whether Terraform will be prevented from destroying the resource. Defaults to "DELETE".
    When a 'terraform destroy' or 'terraform apply' would delete the resource,
    the command will fail if this field is set to "PREVENT" in Terraform state.
    When set to "ABANDON", the command will remove the resource from Terraform
    management without updating or deleting the resource in the API.
    When set to "DELETE", deleting the resource is allowed.
```
{{< /tab >}}
{{% /tabs %}}