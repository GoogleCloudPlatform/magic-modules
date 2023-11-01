# Diff processor

Provides tools that operate on diffs between provider versions.

## Run

```bash
# set up old / new dirs
make clone OWNER_REPO=modular-magician/terraform-provider-google

# build based on old / new dirs
make build OLD_REF=branch_or_commit NEW_REF=branch_or_commit

# Run the commands
bin/diff-processor breaking-changes
GITHUB_TOKEN=github_token bin/diff-processor add-labels PR_ID [--dry-run]  # Requires write access to issues
```

## Test
```bash
go test ./...
```
