# Breaking change detector

## Purpose
Detects breaking changes between provider versions.

Specifically protects customer expectations between [minor version](https://www.terraform.io/plugin/sdkv2/best-practices/versioning#example-minor-number-increments).


## Execution of

### Program:mode-default
```bash
go run .
```

### Tests
```bash
go test ./...
```


## Misc
```bash
# getting the go version label from git log
TZ=UTC0 git log --since="jan 1 2019" --format=%cd-%H --date=format-local:%Y%m%d%H%M%S | sed -E "s/(.*-.{12}).*/\1/"
```


