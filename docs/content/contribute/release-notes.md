---
title: "Write release notes"
weight: 20
---

# Write release notes

This guide explains best practices for composing accurate, end-user focused release notes for Magic Modules pull requests.

Every pull request must have at least one release note block in the opening comment. Release note blocks have the following format:

~~~markdown
```release-note:TYPE
CONTENT
```
~~~

Replace `TYPE` with the correct release note type, and `CONTENT` with a release note written according to the guidelines in the following sections.

## General guidelines

Do | Don't
-- | -----
Only have one `CONTENT` line per release note block. Use multiple blocks if there are multiple related changes in a single PR. | Don't add multiple lines to a single release note block. Avoid combining multiple distinct types of changes into one release block.
If a change only affects the `google-beta` provider add `(beta)` to the end of the release note. If a change only affects the `google` provider add `(ga)` to the end of the release note. | Don't add either suffix if the change affects both providers.
Set an appropriate release note type. | Don't leave the type as `REPLACEME`.

## Type-specific guidelines and examples

{{< tabs >}}
{{< tab "New field(s)" >}}
Write your release note in the following format:

~~~markdown
```release-note:enhancement
PRODUCT: added `FIELD_1`, `FIELD_2`, and `FIELD_N` fields to `RESOURCE_NAME` resource
```
~~~

Replace `PRODUCT`, `FIELD_*`, and `RESOURCE_NAME` according to the pull request content. For example:

~~~markdown
```release-note:enhancement
compute: added `foo_bar` field to `google_compute_foo` resource
```
~~~
{{< /tab >}}
{{< tab "New resource" >}}
Write your release note in the following format:

~~~markdown
```release-note:new-resource
`RESOURCE_NAME`
```
~~~

Replace `RESOURCE_NAME` according to the pull request content. For example:

~~~markdown
```release-note:new-resource
`google_compute_new_resource`
```
~~~
{{< /tab >}}
{{< tab "New datasource" >}}
Write your release note in the following format:

~~~markdown
```release-note:new-datasource
`DATASOURCE_NAME`
```
~~~

Replace `DATASOURCE_NAME` according to the pull request content. For example:

~~~markdown
```release-note:new-datasource
`google_compute_new_datasource`
```
~~~
{{< /tab >}}
{{< tab "Other" >}}
### Choose a release note type
For each release note block, choose an appropriate type from the following list:

- `enhancement` : New features on existing resources
- `bug` : Bug fix
- `deprecation` : A field/resource is being marked as deprecated (not being removed)
- `breaking-change` : Changes that require users to change their configuration
- `note` : General type for other notes that might be relevant to users but don't fit into another category
- `none` : Changes where there is no user impact, like test fixes, website updates and
  CI changes. Release notes of this type should be empty.

### Guidelines

Do | Don't
-- | -----
Use past tense to describe the end state after the change is released. Start with a verb. For example, "added...", "fixed...", or "resolved...". You can use future tense to describe future changes, such as saying that a deprecated field will be removed in a future version. | Don't use present or future tense to describe changes that are included in the pull request.
Write user-focused release notes. For example, reference specific impacted terraform resource and field names, and discuss changes in behavior users will experience. | Avoid API field/resource/feature names. Avoid implementation details. Avoid language that requires understanding of provider internals.
Surround resource or field names with backticks. | Don't use resource or field names without punctuation or with other punctuation like quotation marks.
Use impersonal third person. | Don't use "I", "you", etc.
If the pull request impacts any specific, begin your release note with that product name followed by a colon. Use lower case for the first letter after the colon. For example, `cloudrun: added...` For MMv1 resources, use the folder name that contains the yaml files as the product name; for handwritten or tpgtools resources, use the API subdomain; for broad cross-product changes, use `provider`. | Don't begin your release note with the full resource name. Don't add backticks around the product name. Don't capitalize the first letter after the colon.

### Examples

~~~markdown
```release-note:bug
cloudrun: fixed perma-diff in `google_cloud_run_service`
```
~~~

~~~markdown
```release-note:deprecation
container: deprecated `region` and `zone` on `google_container_unicorn`. Use `location` instead.
```
~~~
{{< /tab >}}
{{< /tabs >}}
