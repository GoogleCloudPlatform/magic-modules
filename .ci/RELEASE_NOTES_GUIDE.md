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

## Headings

Release notes should be formatted with one of the following headings:
- `release-note:enhancement`
- `release-note:bug`
- `release-note:note`
- `release-note:new-resource`
- `release-note:new-datasource`
- `release-note:deprecation`
- `release-note:breaking-change`
- `release-note:none`

However, any note with a language heading starting with ```release-note:... will get copied.

## Non-User-Facing PRs

Any PR that should not have any impact on users (test fixes, code generation, website updates,
CI changes, etc.) should use a `release-note:none` block. It can be left empty, or can be
optionally filled with an explanation of why the PR is not user-facing.

By including this block explicitly, it lets whoever is generating the changelog know that
a release note was explicitly omitted, not just forgotten. It'll also let your PR pass any
future automation around release note correctness checking.

## Release Note Style Guide (Terraform-specific)

Notes SHOULD:
- Start with a verb
- Use past tense (added/fixed/resolved) as much as possible
- Only use present tense to suggest future behavior, i.e. for breaking
  changes, deprecations, or new behavior.
- Impersonal third person (no “I”, “you”, etc.)
- Start with {{service}} if changing an existing resource (see below)

Notes, breaking changes, and features are exceptions. These are more free-form and left to
the discretion of the PR author and reviewer. The overarching goal should be a good user
experience when reading the changelog.

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

**Notes:** General tag for things that don’t have changes in provider but may be important to users. Syntax is slightly more flexible here. 

```release-note:note
Starting on Nov 1, 2019, Cloud Functions API will be private by default. Add appropriate bindings through `google_cloud_function_iam_*` resources to manage privileges for `google_cloud_function` resources created by Terraform.
```
### Counter-examples:

The following changelog entries are not ideal.

#### No Type

~~~
```release-note:REPLACEME
compute: fixed permadiff on description for `google_compute_instance`
```
~~~

This doesn't update the type of release note, which means it will need to be corrected at generation time.

Better:

~~~
```release-note:bug
compute: fixed permadiff on description for `google_compute_instance`
```
~~~

### Not Past Tense

~~~
```release-note:bug
compute: Fix permadiff on description for `google_compute_instance`
```
~~~

This doesn't use the past tense. Readers of teh changelog will be reading what _happened_ in a release,
so the language should be that of describing what happened. Imagine you're answering the question
"what changed since the last version?"

Better:

~~~
```release-note:bug
compute: Fixed permadiff for `google_compute_instance`
```
~~~

### No Service

~~~
```release-note:bug
Fixed permadiff on description for `google_compute_instance`
```
~~~

This doesn't start with a service name. By convention, we prefix all our bug and enhancement changelog
entries with service names, and the other entries when it makes sense and seems beneficial. This helps
sort the changelog and group related changes together, and helps users scan for the services they use.

Better:

~~~
```release-note:bug
compute: Fixed permadiff for `google_compute_instance`
```
~~~

### Not User-Centric

~~~
```release-note:bug
compute: made description Computed for `google_compute_instance`
```
~~~

This isn't written for the right audience; our users don't all, or even mostly, know what Computed
means, and shouldn't have to. Instead, describe the impact that this will have on them.

Better:

~~~
```release-note:bug
compute: fixed permadiff on description for `google_compute_instance`
```
~~~

### Resource Instead of Service

~~~
```release-note:bug
compute_instance: Fixed permadiff on description for `google_compute_instance`
```
~~~

This uses the resource instead of the service as a prefix. 

### Unticked Resource Names

~~~
```release-note:bug
compute: Fixed permadiff on description for google_compute_instance
```
~~~

This doesn't have \`\` marks around the resource name, which by convention we do. This sets the resource
name apart, making it easer to notice.

Better:

~~~
```release-note:bug
compute: Fixed permadiff on description for `google_compute_instance`
```
~~~

Choosing the right service name is a bit of an art. A good rule of thumb is if there's something
besides the resource name after `google_`, use that. For example, `compute` is a good choice from
`google_compute_instance`. Not every resource has that, however; for `google_project`, the service
is not part of the resource address. In these cases, falling back on the name of the package the
resource's APIs is implemented in (`resourcemanager`, for `google_project`) is a good call.

Not every change applies only to one resource. Judgment is best here. When in doubt, `provider` is
a good way to indicate sweeping changes that are likely to impact most users.
