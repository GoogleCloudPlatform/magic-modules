package metadata

import (
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
)

func FromProperties(props []*api.Type) []Field {
	var fields []Field
	for _, p := range props {
		f := Field{
			Json:         p.IsJsonField(),
			ProviderOnly: p.ProviderOnly(),
		}
		if !p.ProviderOnly() {
			f.ApiField = p.MetadataApiLineage()
		}
		if p.ProviderOnly() || p.MetadataLineage() != p.MetadataDefaultLineage() {
			f.Field = p.MetadataLineage()
		}
		fields = append(fields, f)
	}
	return fields
}

// Field is a field in a metadata.yaml file.
type Field struct {
	// The name of the field in the REST API, including the path. For example, "buildConfig.source.storageSource.bucket".
	ApiField string `yaml:"api_field,omitempty"`
	// The name of the field in Terraform, including the path. For example, "build_config.source.storage_source.bucket". Defaults to the value
	// of `api_field` converted to snake_case.
	Field string `yaml:"field,omitempty"`
	// If true, the field is only present in the provider. This primarily applies for virtual fields and url-only parameters. When set to true,
	// `field` should be set and `api_field` should be left empty. Default: `false`.
	ProviderOnly bool `yaml:"provider_only,omitempty"`
	// If true, this is a JSON field which "covers" all child API fields. As a special case, JSON fields which cover an entire resource can
	// have `api_field` set to `*`.
	Json bool `yaml:"json,omitempty"`
}
