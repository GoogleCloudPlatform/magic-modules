# Breaking change detector

Detects breaking changes between provider versions.

Specifically protects customer expectations between [minor version](https://www.terraform.io/plugin/sdkv2/best-practices/versioning#example-minor-number-increments).


## Run

```bash
# set up old / new dirs
make clone OWNER_REPO=modular-magician/terraform-provider-google

# build based on old / new dirs
make build OLD_REF=branch_or_commit NEW_REF=branch_or_commit

# Run the binary
bin/breaking-change-detector
```

## Test
```bash
go test ./...
```
