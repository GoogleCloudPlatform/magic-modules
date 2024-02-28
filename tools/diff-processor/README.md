# Diff processor

Provides tools that operate on diffs between provider versions.

## Run

```bash
# set up old / new dirs
make clone OWNER_REPO=modular-magician/terraform-provider-google

# build based on old / new dirs
make build OLD_REF=branch_or_commit NEW_REF=branch_or_commit

# Run breaking change detection on the difference between OLD_REF and NEW_REF
bin/diff-processor breaking-changes

# Add labels to a PR based on the resources changed between OLD_REF and NEW_REF
# The token used must have write access to issues
GITHUB_TOKEN_MAGIC_MODULES=github_token bin/diff-processor add-labels PR_ID [--dry-run]
```

## Test
```bash
go test ./...
```
