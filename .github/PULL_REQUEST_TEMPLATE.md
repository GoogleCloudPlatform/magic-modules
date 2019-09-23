<!--
Note: You may see "This branch is out-of-date with the base branch"
when you submit a pull request. This is fine! We don't use the GitHub
merge button to merge PRs, and you can safely ignore that message.

Thanks for contributing!
-->

<!-- AUTOCHANGELOG for Downstream PRs.
EXTERNAL CONTRIBUTORS: Your reviewer will most likely fill this in for you, so don't worry about this section!

For some repos (currently Terraform GA/beta providers), we have the
ability to autogenerate CHANGELOGs.

NO CHANGELOG NOTE: If you do not want a release note,
please add the "changelog: no-release-note" label to this PR.

Otherwise, fill the template out below
-->

# Release Note Template for Downstream PRs (will be copied)
```release-note:enhancement

```

<!-- GUIDE FOR WRITING RELEASE NOTES
Release notes should be formatted with one of the following headings.
- release-note:bug
- release-note:note
- release-note:new-resource
- release-note:new-datasource
- release-note:deprecation
- release-note:breaking-change

Guide for writing release notes:

Notes SHOULD:
- Start with a verb
- Use past tense (added/fixed/resolved) as much as possible
- Only use present tense in imperative sentences to suggest future behavior for
  breaking changes/deprecations ("Use X" vs "You should use X" or "Users should use X")
- Impersonal third person (no “I”, “you”, etc.)
- Start with `{{service}}` if changing an existing resource (see exampels below)

DO:

```release-note:enhancement
compute: added `foo_bar` field to `google_compute_foo` resource
```

```release-note:bug
container: fixed perma-diff in `google_container_cluster`
```

```release-note:breaking-change
project: made `iam_policy` authoritative
```

```release-note:deprecation
container: deprecated `region` and `zone` on `google_container_unicorn`. Use `location` instead.
```

Note no service name or *New Resource* tag:
```release-note:new-resource
`google_compute_new_resource`
```

Note no service name or *New Datasource* tag:
```release-note:new-datasource
`google_compute_new_datasource`
```

DON'T DO:
- Add compute_instance resource
- Fix bug
- fixed a bug in google_compute_network
- `google_project` now supports `blah`
- You can now create google_sql_instances in us-central1
- Adds support for `google_source_repo_repository`’s `url` field
- Users should now use location instead of zone/region on `google_container_unicorn`
-->
