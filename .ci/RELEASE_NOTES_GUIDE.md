# Release Notes for Terraform/Magic Modules Auto-CHANGELOG

## Background

The Magician bot has the ability to copy specifically formatted release notes
from upstream Magic-Modules PRs to downstream PRs that have CHANGELOGS, namely
PRs generated in the Terraform providers (GA, beta).

Code lives in magic-modules/downstream_changelog_metadata.py

This guide discusses the style and format of release notes to add
in PR descriptions so they will be copied downstream and used
in CHANGELOGs.

## Expected Format

The description should have Markdown-formatted code blocks with language
headings (i.e. text right after the three ticks) like this:

~~~
PR Description

...

```release-note:enhancement
compute: Fixed permadiff for `bar` in `google_compute_foo`
```
~~~

You can have multiple code blocks to have multiple release notes per PR, i.e.

~~~

PR Description

...

```release-note:deprecation
container: Deprecated `region` for `google_container_cluster` - use location instead.
```

```release-note:enhancement
container: Added general field `location` to `google_container_cluster`
```
~~~


Do not indent the block and make sure to leave newlines, so you don't confuse
the Markdown parser.


## Release Note Style Guide (Terraform-specific)

Notes SHOULD:
- Start with a verb
- Use past tense (added/fixed/resolved) as much as possible
- Only use present tense to suggest future behavior, i.e. for breaking
  changes, deprecations, or new behavior.
- Impersonal third person (no “I”, “you”, etc.)
- Start with {{service}} if changing an existing resource (see below)

See examples below for good release notes.

### Examples:

**Enhancements:** adding fields or new features to existing resources

~~~
```release-note:enhancement
compute: added `foo_bar` field to `google_compute_foo` resource
```
~~~
**Bugs:** fixing existing resources

~~~
```release-note:bug
container: fixed perma-diff in `google_container_cluster`
```
~~~

**Breaking Changes:** changes to existing resources that may require users to change their config

~~~
```release-note:breaking-change
project: made `iam_policy` authoritative
```
~~~

**Deprecations:** Announce deprecations when fields/resources are marked as deprecated, not removed

~~~
``` release-note:deprecation
container: deprecated `region` and `zone` on `google_container_unicorn`. Use `location` instead.
```
~~~

**New Resources And Datasources:**
(note no service name or *New Resource* tag)

~~~
```release-note:new-resource
`google_compute_new_resource`
```
~~~

~~~
```release-note:new-datasource
`google_compute_new_datasource`
```
~~~

Notes: General tag for things that don’t have changes in provider but may be important to users. Syntax is slightly more flexible here. 

```release-note:note
Starting on Nov 1, 2019, Cloud Functions API will be private by default. Add appropriate bindings through `google_cloud_function_iam_*` resources to manage privileges for `google_cloud_function` resources created by Terraform.
```

Don’t write notes like:
- Add compute_instance resource
  - not past tense
  - no \`\` around resource
  - doesn't start with service
- Fix bug
  - not past tense
  - no indication of impact to users
  - doesn't start with service
- fixed a bug in google_compute_network
  - no \`\` around resource
  - doesn't start with service
  - unclear impact to users
- `google_project` now supports `blah`
  - not past tense
  - doesn't start with service
- You can now create google_sql_instances in us-central1
  - not past tense
  - second person voice
  - no \`\` around resource
  - doesn't start with service
- Adds support for `google_source_repo_repository`’s `url` field
  - not past tense
  - doesn't start with service
- Users should now use location instead of zone/region on `google_container_unicorn`
  - prescriptive instead of descriptive
  - not past tense
  - doesn't start with service
  - would be fine after a "container: deprecated `location` field on `google_container_unicorn`."

## Headings

Release notes should be formatted with one of the following headings:
- `release-note:enhancement`
- `release-note:bug`
- `release-note:note`
- `release-note:new-resource`
- `release-note:new-datasource`
- `release-note:deprecation`
- `release-note:breaking-change`

However, any note with a language heading starting with ```release-note:... will get copied.
