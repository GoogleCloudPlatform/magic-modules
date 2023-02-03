# Release Notes for Terraform/Magic Modules Auto-CHANGELOG

## Background

The Magician bot has the ability to copy specifically formatted release notes
from upstream Magic-Modules PRs to downstream PRs that have CHANGELOGS, namely
PRs generated in the Terraform providers (GA, beta).

This guide discusses the style and format of release notes to add
in PR descriptions so they will be copied downstream and used
in CHANGELOGs ([GA provider changelog](https://github.com/hashicorp/terraform-provider-google/blob/main/CHANGELOG.md), 
[beta provider changelog](https://github.com/hashicorp/terraform-provider-google-beta/blob/main/CHANGELOG.md)).

## Expected Format

The description should have Markdown-formatted code blocks with language
headings (i.e. text right after the three ticks) like this:

~~~
<PR Description

...>

```release-note:enhancement  <-- change type
compute: added `foo_bar` field to `google_compute_foo` resource <-- change description
```
~~~

All pull requests need to include release-note block and **change type** (`release-note:REPLACEME` in the PR template) should be replaced with one of the following [headings](#headings), while **change description** can be omitted from some PRs with non-user-facing changes

### Headings

Release notes should be formatted with one of the following headings (See [examples](#examples) below for good release notes for each change type):
| Headings                       | Usage                                                                               |
| ------------------------------ | ----------------------------------------------------------------------------------- |
| `release-note:enhancement`     | adding fields or new features to existing resources                                 |
| `release-note:bug`             | fixing existing resources                                                           |
| `release-note:note`            | things that don’t have changes in provider but may be important to users            |
| `release-note:new-resource`    | adding new resources                                                                |
| `release-note:new-datasource`  | adding new datasources                                                              |
| `release-note:deprecation`     | announcing deprecations when fields/resources are marked as deprecated, not removed |
| `release-note:breaking-change` | changes to existing resources that may require users to change their config         | 
| `release-note:none`            | changes that don't have any impact on users                                         | 

### Release Note Style Guide (Terraform-specific)

Notes SHOULD:
- Start with a verb
- Use past tense (added/fixed/resolved) as much as possible
- Only use present tense to suggest future behavior, i.e. for breaking
  changes, deprecations, or new behavior.
- Not be capitalized            
- Impersonal third person (no “I”, “you”, etc.)
- Start with {{service}} if changing an existing resource (i.e. for `enhancement` and `bug`)
- Include provider resource name (with `google_` prefix) if changing an existing resource
- Use \`\` marks around the resource and field names
- Use user-centric language for description

Notes, breaking changes are exceptions. These are more free-form and left to the discretion of the PR author and reviewer. The overarching goal should be a good user experience when reading the changelog.

### Examples:
* `release-note:enhancement`
~~~
```release-note:enhancement
compute: added `foo_bar` field to `google_compute_foo`
```
~~~

* `release-note:bug`
~~~
```release-note:bug
compute: fixed perma-diff on `region` field in `google_compute_foo` resource
```
~~~

* `release-note:note`
~~~
```release-note:note
updated Bigtable go client version from 1.13 to 1.16
```
~~~

* `release-note:new-resource`
~~~
```release-note:new-resource
`google_compute_new_resource`
```
~~~

* `release-note:new-datasource`
~~~
```release-note:new-resource
`google_compute_new_datasource`
```
~~~

* `release-note:deprecation`
~~~
```release-note:deprecation
compute: deprecated `region` field for `google_compute_foo` - use location instead.
```
~~~

* `release-note:breaking-change`
~~~
```release-note:breaking-change
compute: changeed `region` field for `google_compute_foo` to be required
```
~~~

* `release-note:none`  
~~~
```release-note:none
(change description can be left empty, or can be optionally filled with an explanation of why the PR is not user-facing)
```
~~~

### Beta-only changes

To qualify that a change is specific to the beta provider add `(beta)`
at the end of the release note.
This will notify whoever is generating the release notes omit the note from changelogs for the ga release.

~~~
```release-note:enhancement
compute: added `foo_bar` field to `google_compute_foo` (beta)
```
~~~

### Revert changes

Copy the original PR release note block and add `(revert)` at the end.
This will let whoever is generating the release notes omit the original PR note from chagnelogs

~~~
```release-note:enhancement
compute: added `foo_bar` field to `google_compute_foo` (revert)
```
~~~

### Multiple code blocks per PR

You can have multiple code blocks to have multiple release notes per PR, i.e.
~~~
<PR Description

...>

```release-note:enhancement
compute: added `foo_bar` field to `google_compute_foo`
```

```release-note:deprecation
compute: deprecated `region` field for `google_compute_foo` - use location instead.
```
~~~

Do not indent the block and make sure to leave newlines, so you don't confuse
the Markdown parser.

### Counter-examples (Common Mistakes):

The following changelog entries are not ideal.

#### No Type

~~~
```release-note:REPLACEME
compute: fixed permadiff on `description` for `google_compute_instance`
```
~~~

This doesn't update the type of release note, which means it will need to be corrected at generation time.

Better:

~~~
```release-note:bug
compute: fixed permadiff on `description` for `google_compute_instance`
```
~~~

#### Not Past Tense

~~~
```release-note:bug
compute: fix permadiff on description for `google_compute_instance`
```
~~~

This doesn't use the past tense. Readers of the changelog will be reading what _happened_ in a release,
so the language should be that of describing what happened. Imagine you're answering the question
"what changed since the last version?"

Better:

~~~
```release-note:bug
compute: fixed permadiff on `description` for `google_compute_instance`
```
~~~

#### No Service

~~~
```release-note:bug
fixed permadiff on description for `google_compute_instance`
```
~~~

This doesn't start with a service name. By convention, we prefix all our bug and enhancement changelog
entries with service names, and the other entries when it makes sense and seems beneficial. This helps
sort the changelog and group related changes together, and helps users scan for the services they use.

Better:

~~~
```release-note:bug
compute: fixed permadiff on `description` for `google_compute_instance`
```
~~~

### Resource Instead of Service

~~~
```release-note:bug
compute_instance: fixed permadiff on `description` for `google_compute_instance`
```
~~~

This uses the resource instead of the service as a prefix. 

Better:

~~~
```release-note:bug
compute: fixed permadiff on `description` for `google_compute_instance`
```
~~~

Choosing the right service name is a bit of an art. A good rule of thumb is if there's something
besides the resource name after `google_`, use that. For example, `compute` is a good choice from
`google_compute_instance`. Not every resource has that, however; for `google_project`, the service
is not part of the resource address. In these cases, falling back on the name of the package the
resource's APIs is implemented in (`resourcemanager`, for `google_project`) is a good call.

Not every change applies only to one resource. Judgment is best here. When in doubt, `provider` is
a good way to indicate sweeping changes that are likely to impact most users.

#### Not User-Centric

~~~
```release-note:bug
compute: made `description` Computed for `google_compute_instance`
```
~~~

This isn't written for the right audience; our users don't all, or even mostly, know what Computed
means, and shouldn't have to. Instead, describe the impact that this will have on them.

Better:

~~~
```release-note:bug
compute: fixed permadiff on `description` for `google_compute_instance`
```
~~~

#### Unticked Resource/ Field Names

~~~
```release-note:bug
compute: fixed permadiff on description for google_compute_instance
```
~~~

This doesn't have \`\` marks around the resource name `google_compute_instance` and field name `description`, which by convention we do. This sets the resource
name and field name apart, making them easer to notice.

Better:

~~~
```release-note:bug
compute: fixed permadiff on `description` for `google_compute_instance`
```
~~~

