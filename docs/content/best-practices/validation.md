---
title: "Validation"
weight: 50
---

# Validation

There are a number of ways to add client-side validation to resources. The benefit of client-side validation is that errors can be surfaced at plan time, instead of partway through a (potentially very long) apply process, allowing for faster iteration. However, the tradeoff is that client-side validation can get out of sync with server-side validation, creating additional maintenance burden for the provider and preventing users from accessing the latest features without upgrading.

Client-side validation is generally discouraged due to the low positive impact of an individual validation rule and outsized negative impact when client-side validation and API capabilities drift, requiring both provider changes and users to update. Client-side validation may be added in cases where it is extremely unlikely to change, covered below.

In theory, APIs that have a validation endpoint could use it to ensure that client-side validation always matches server-side validation. However, this is not a well-lit path. Follow [hashicorp/terraform-provider-google#20713](https://github.com/hashicorp/terraform-provider-google/issues/20713) for more information.

The following sections cover best practices for specific types of client-side validation.

## URL segments

If a resource URL looks like:

```
projects/{project}/folders/{folder}/resource/{resource_id}
```

Adding validation for the last part of the path (`resource_id`) is likely safe in most cases, especially for GCE resources.

## Enum

Enums are generally okay if they are exhaustive of all possible values for a clearly defined domain, i.e. where new values are extremely unlikely. Otherwise, it is better to use a string field and add a link to the API documentation as a reference for the possible values.

## Inter-field relationships

`conflicts_with`, `required_with`, and similar are safe types of client-side validation, because they are intrinsically linked to fields in the provider. If a new field is added to the API that invalidates a rule, users will need to update the provider to get access to that field anyway, so there isn't a future-compatibility concern.

## Immutable facts

It is safe to validate things that will definitely always be true about an API. For example, a `node_count` field will most likely always need to be non-negative. That is safe to validate. However, validating a max value for `node_count` may not be safe, because the API might increase the allowed values in the future.
