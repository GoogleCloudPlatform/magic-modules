# changelog-pr-body-check

`changelog-pr-body-check` is a command that will ensure that the body of a PR
has at least one changelog entry in it, and that all changelog entries are
valid. If no changelog entry is found, or one or more invalid changelog entries
are found, `changelog-pr-body-check` will comment on the PR to inform the
author of the issue.

Right now, acceptance criteria are hardcoded simply as being one of the
following types of entries:

* bug
* note
* enhancement
* new-resource
* new-datasource
* deprecation
* breaking-change
* feature

A configuration system is planned to allow a more customizable check.

## Usage

This binary requires three environment variables to be set:

* `GITHUB_REPO`, the GitHub repository the PR being checked lives in.
* `GITHUB_OWNER`, the owner of the GitHub repository the PR being checked lives
  in.
* `GITHUB_TOKEN`, an access token with permission to read and comment on issues
  in the repository the PR being checked lives in.

Once these environment variables are set, run the command:

```sh
$ changelog-pr-body-check $NUMBER
```

where `NUMBER` is the ID of the PR to check.

## Results

Any failures will be logged to stderr. If the check passes, it will return
status code 0. Status code 1 indicates that either the PR did not pass all the
checks, and should have comments on it and log entries explaining what failed,
or there was an error running the checks or leaving comments, and stderr should
have more details on what went wrong.
