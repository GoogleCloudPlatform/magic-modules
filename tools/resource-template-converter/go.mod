module github.com/GoogleCloudPlatform/magic-modules/tools/resource-template-converter

go 1.26.0

require github.com/GoogleCloudPlatform/magic-modules/mmv1 v0.0.0

replace github.com/GoogleCloudPlatform/magic-modules/mmv1 => ../../mmv1

require (
	github.com/spf13/cobra v1.10.1
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
)
