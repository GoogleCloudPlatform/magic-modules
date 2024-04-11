---
subcategory: "Cloud Platform"
description: |-
 Allows removal of a member:role pairing from an IAM policy.
---

# google\_project\_iam\member\_remove

Allows removal of a member:role pairing from an IAM policy. For more information see
[the official documentation](https://cloud.google.com/iam/docs/granting-changing-revoking-access)
and
[API](https://cloud.google.com/resource-manager/reference/rest/v1/projects/setIamPolicy).


```hcl
resource "google_project_iam_member_remove" "foo" {
  role     = "roles/editor"
  project  = "your-project-id"
  member  = "default-gce-sa@example.com"
}
```

## Argument Reference

The following arguments are supported:

* `member` - (Required) Identities that will be granted the privilege in `role`.
  Each entry can have one of the following values:
  * **user:{emailid}**: An email address that represents a specific Google account. For example, alice@gmail.com or joe@example.com.
  * **serviceAccount:{emailid}**: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com.
  * **group:{emailid}**: An email address that represents a Google group. For example, admins@example.com.
  * **domain:{domain}**: A G Suite domain (primary, instead of alias) name that represents all the users of that domain. For example, google.com or example.com.

* `role` - (Required) The role that should be removed. 

* `project` - (Required) The project id of the target project.
