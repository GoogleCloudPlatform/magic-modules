---
title: "Add handwritten IAM resources"
weight: 23
---

# Add handwritten IAM resources

Handwritten IAM support is only recommended for resources that cannot be managed
using [MMv1](/magic-modules/docs/how-to/add-mmv1-iam),
including for handwritten resources, due to the need to manage tests and
documentation by hand. This guidance goes through the motions of adding support
for new handwritten IAM resources, but does not go into the details of the
implementation as any new handwritten IAM resources are expected to be
exceptional.

IAM resources are implemented using an IAM framework, where you implement an
interface for each parent resource supporting `getIamPolicy`/`setIamPolicy` and
the associated IAM resources that target that parent resource- `_member`,
`_binding`, and `_policy`- are created by the framework.

To add support for a new target, create a new file in
`mmv1/third_party/terraform/utils` called `iam_{{resource}}.go`, and implement
the `ResourceIamUpdater`, `newResourceIamUpdaterFunc`, `iamPolicyModifyFunc`,
`resourceIdParserFunc` interfaces from
https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/utils/iam.go.erb
in public types, alongside a public `map[string]*schema.Schema` containing all
fields referenced in the resource.

Once your implementation is complete, add the IAM resources to `provider.go`
inside the `START non-generated IAM resources` block, creating the concrete
resource types using the `ResourceIamMember`, `ResourceIamBinding`, and
`ResourceIamPolicy` functions. For example:

```go
				"google_bigtable_instance_iam_binding":         ResourceIamBinding(IamBigtableInstanceSchema, NewBigtableInstanceUpdater, BigtableInstanceIdParseFunc),
				"google_bigtable_instance_iam_member":          ResourceIamMember(IamBigtableInstanceSchema, NewBigtableInstanceUpdater, BigtableInstanceIdParseFunc),
				"google_bigtable_instance_iam_policy":          ResourceIamPolicy(IamBigtableInstanceSchema, NewBigtableInstanceUpdater, BigtableInstanceIdParseFunc),
```

Following that, write a test for each resource exercising create and update for
both `_policy` and `_binding`, and create for `_member`. No special
accommodations are needed for the IAM test compared to a normal Terraform
resource test.

Documentation for IAM resources is done using single page per target resource,
rather than a distinct page for each IAM resource level. As most of the page is
standard, you can generally copy and edit an existing handwritten page such as
https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/website/docs/r/bigtable_instance_iam.html.markdown
to write the documentation.
