module github.com/GoogleCloudPlatform/magic-modules/tools/sweeper-resolver

go 1.23.3

replace github.com/GoogleCloudPlatform/magic-modules/tools/test-reader => ../test-reader

replace github.com/GoogleCloudPlatform/magic-modules/mmv1 => ../../mmv1

require (
	github.com/GoogleCloudPlatform/magic-modules/mmv1 v0.0.0-00010101000000-000000000000
	github.com/GoogleCloudPlatform/magic-modules/tools/test-reader v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/agext/levenshtein v1.2.1 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/golang/glog v1.2.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/hashicorp/hcl/v2 v2.20.1 // indirect
	github.com/mitchellh/go-wordwrap v0.0.0-20150314170334-ad45545899c7 // indirect
	github.com/zclconf/go-cty v1.13.0 // indirect
	golang.org/x/exp v0.0.0-20240222234643-814bf88cf225 // indirect
	golang.org/x/mod v0.15.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	golang.org/x/tools v0.18.0 // indirect
)
