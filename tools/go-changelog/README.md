# go-changelog

go-changelog is a library and a set of binaries for working with changelog
generation, written in Go.

The underlying strategy is to have a directory committed as part of the git
repo you're generating a changelog for--by convention, `.changelog` is used,
but the tool doesn't care what it's called or where it is. Inside that
directory, each unit of change--almost always a pull request--has a file named
after the unit of change. For example, a PR would include a `.changelog/PR#.txt`
file, where `PR#` is the unique ID for the PR.

Because each PR or unit of change has its own file, you should never encounter
merge conflicts for changelog entries.

The tool then takes two commits for the repository--specified by any valid git
reference--and compiles a list of the files present in the later commit that
aren't present in the earlier commit. These are assumed to accurately reflect
the changes that have been added between the two commits.

## Installation

### Binaries

To run the command line binaries, use Go to build the commands. For example:

```sh
$ go install github.com/hashicorp/go-changelog/cmd/changelog-pr-body-check@latest
```

### Docker

A Dockerfile is provided that will build an image containing the binaries. You
can either run this container directly, or you can build an image sourced from
it that specifies the environment variables and entrypoint. In the future, when
go-changelog has config files, this will be how you can add your config files.

```Dockerfile
FROM hashicorpdev/go-changelog

ENV GITHUB_REPO=myrepo
ENV GITHUB_OWNER=myorg

# Maybe leave this one out and specify it with -e
ENV GITHUB_TOKEN=foo123abc

ENTRYPOINT ["/go-changelog/changelog-pr-body-check"]
```

## Usage

For using the `go-changelog` library, please see [go.dev](https://pkg.go.dev/github.com/hashicorp/go-changelog).
For using the specific binaries, please see the README files in their
directories.

## Change File Formatting

The files in your directory describe the changes that will be used to generate
the changelog. Their contents should have the following formatting:

~~~
```release-note:TYPE
ENTRY
```
~~~

Where `TYPE` is the type of release note entry this is. This is usually "bug",
"enhancement", etc. The tool does not prescribe a list of types to choose from;
whatever you enter will be available to you when generating the changelog.

`ENTRY` is the body of the changelog entry, and should describe the changes
that were made. This is used as free-text input and will be returned to you as
it is entered when generating the changelog.

Sometimes PRs have multiple changelog entries associated with them. In this
case, use multiple blocks.

~~~
```release-note:deprecation
Deprecated the `foo` interface, please use the `bar` interface instead.
```

```release-note:enhancement
Added the `bar` interface.
```
~~~

## Best Practices

### Keep changelog entries with code change commits

This system works under the assumption that changelog entries will always be
kept with the code they're describing changes to. So it works best if the code
and its changelog entry appear in the same commit. This means cherry-picks that
move the code around will make sure to bring the changelog entries with it and
will keep the changelog correct.

Using squash merges makes this easier.

### Lean on automation

A sample `changelog-check` binary is included in the `cmd` directory to show
how a GitHub PR can be checked to ensure that a changelog entry is attached to
the PR. Lean on automation to guard against forgetting changelog entries when
submitting PRs.

You can also have bots generate these files from PR bodies--when a PR body is
updated, have a bot push a commit updating the changelog entry as well. The
markdown code blocks can be easily parsed out of the PR body using the
`NotesFromEntry` function in the `changelog` package, and re-formatted into a
file for the commit. This means users don't need to mess with git to update
changelog entries.

## Shortcomings

### Immutable changelog entries

Once a changelog entry is merged, it's hard to mutate it and make sure the
commit updating it stays with the code. It's not impossible--you just need
to remember to always cherrypick _both_ commits--but it is harder. It is
recommended that you use automation and PR review to make sure you're happy
with the changelog entry before merging, as much as possible, to keep these
situations limited.

### Silent failure when no changelog entry is set

If a change forgets to include a changelog entry at all, the tool will ignore
it completely. PR-based tools can shout loudly that a PR did not include a
change, but this system does not have that capability without expecting every
commit to have a changelog entry associated with it, which seems unreasonable.

It is recommended to use automation to prevent changes without changelog
entries before they're merged, to ensure every change gets an entry.

## Why

Changelogs are an important interface for helping users be aware of development
efforts and behavioral changes in software. Writing good changelogs is tricky.

* Writing changelogs by hand inevitably leads to people forgetting to include
  things in the changelog.
* Basing changelog entries on commit messages conflates two audiences; commit
  messages are for developers and maintainers and should be directed towards
  them, changelog entries are directed towards users and should be written in
  the terminology they understand and use.
* Writing changelogs in PR bodies is quick, easy, and allows them to be part
  of the PR review process, but PR bodies aren't stored in the repository
  itself (meaning they can't be easily backed up, etc.) and makes it harder
  to keep changelogs with the code they describe as that code is cherry-picked,
  backported, and merged across branches.

`go-changelog` prioritises making it easy to generate correct changelogs no
matter what shenanigans you engaged in to build the branch you're releasing
from. It also prioritises flexibility around workflows and changelog
formatting.

## Prior Art

This package is based on a bunch of experiments with the [Terraform provider for Google Cloud](https://github.com/terraform-providers/terraform-provider-google)
and the lessons learned while generating it. It is also based on prior art in
the community:

* [Paul Tyng's changelog-gen](https://github.com/paultyng/changelog-gen)
* [Amber Brown's Towncrier](https://github.com/hawkowl/towncrier)
