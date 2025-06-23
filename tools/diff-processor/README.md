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

# Compute service labels to add bsaed on the resources changed between OLD_REF and NEW_REF
bin/diff-processor changed-schema-labels
```

## Test
```bash
go test ./...
```
