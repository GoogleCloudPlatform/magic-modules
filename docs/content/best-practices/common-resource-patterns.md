---
title: "Common resource patterns"
weight: 40
---

# Common resource patterns

## Singletons

Singletons are resources – often config or settings objects – that can only exist once. In some cases, it may be possible to create and delete the resource (but only one can exist at a time); in other cases the resource _always_ exists and can only be read and updated.

Implementing resources like this may require some or all of the following:

1. If there _isn't_ a create endpoint, set the [create_url]({{< ref "/reference/resource/#create_url" >}}) to point to the update endpoint.
1. If there _is_ a create endpoint, add [pre-create custom code]({{< ref "/develop/custom-code/#pre_post_injection" >}}) that implements "acquire-on-create" logic. The custom code should check whether the resource already exists with a read request, and if it does, run the update logic and return early. For example, see [mmv1/templates/terraform/pre_create/firebasehosting_site.go.tmpl](https://github.com/GoogleCloudPlatform/magic-modules/blob/dc4d9755cb9288177e0996c1c3b3fa9738ebdf89/mmv1/templates/terraform/pre_create/firebasehosting_site.go.tmpl).
   * Note: The main disadvantage of "acquire-on-create" logic is that users will not be presented with a diff between the resource's old and new states – because from the terraform perspective, the resource is only being created. Please upvote https://github.com/hashicorp/terraform/issues/19017 to request better support for this workflow.
1. If there is no delete endpoint, set [`exclude_delete: true`]({{< ref "/reference/resource/#create_url" >}}) at the top level of the resource.

Tests for singletons can run into issues because they are modifying a shared state. To avoid the problems this can cause, ensure that the tests [create dedicated parent resources]({{< ref "/test/test#create-test-projects" >}}) instead of modifying the default test environment. If there need to be multiple test cases, make sure they either have individual parent resources, or that they run serially, like [TestAccAccessContextManager](https://github.com/hashicorp/terraform-provider-google-beta/blob/88fa0756f2ce116765edd4c1551680d9029621f6/google-beta/services/accesscontextmanager/resource_access_context_manager_access_policy_test.go#L31-L33).
