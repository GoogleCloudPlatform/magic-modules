---
title: "Validation"
weight: 50
---

# Validation

There are a number of ways to add client-side validation to resources. The benefit of client-side validation is that errors can be surfaced at plan time, instead of partway through a (potentially very long) apply process, allowing for faster iteration. However, the tradeoff is that client-side validation can get out of sync with server-side validation, creating additional maintenance burden for the provider and preventing users from accessing the latest features without upgrading.

Client-side validation is generally discouraged due to the low positive impact of an individual validation rule and outsized negative impact when client-side validation and API capabilities drift, requiring both provider changes and users to update. Client-side validation may be added in cases where it is extremely unlikely to change, covered below.

The following sections cover best practices for specific types of client-side validation.

## URL segments

If a resource URL looks like:

```
projects/{project}/folders/{folder}/resource/{resource_id}
```

Adding validation for the last part of the path (`resource_id`) may be safe if there are specific restrictions that aren't going to change, such as following an external RFC or other spec/standard. However, if the API was ever less restrictive (or becomes less restrictive later), resources created with other tools and then imported into Terraform may be impossible to actually manage with Terraform (without deleting & recreating them) because the ID which was valid in the API violates the more restrictive validation in the provider.

## Enum

Enums are generally okay if they are exhaustive of all possible values for a clearly defined domain where new values are extremely unlikely. Otherwise, it is better to use a string field and add a link to the API documentation as a reference for the possible values.

## Inter-field relationships

[`conflicts`]({{< ref "/reference/field#conflicts" >}}), [`required_with`]({{< ref "/reference/field#required_with" >}}), [`exactly_one_of`]({{< ref "/reference/field#exactly_one_of" >}}), and [`at_least_one_of`]({{< ref "/reference/field#at_least_one_of" >}}) are often safe to add. However, if there is a chance that the API validation will relax in the future (such as two fields no longer being required together, or two fields no longer conflicting) it's better to not add the restriction in the first place.

## Immutable facts

It is safe to validate things that will definitely always be true about an API. For example, a `node_count` field will most likely always need to be non-negative. That is safe to validate. However, validating a max value for `node_count` may not be safe, because the API might increase the allowed values in the future.
