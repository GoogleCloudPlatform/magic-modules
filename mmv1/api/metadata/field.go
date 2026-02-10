package metadata

import (
	"slices"
	"strings"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
)

func FromProperties(props []*api.Type) []Field {
	// Sort props by lineage
	slices.SortFunc(props, func(a, b *api.Type) int {
		if strings.Join(a.Lineage(), ".") < strings.Join(b.Lineage(), ".") {
			return -1
		}
		return 1
	})

	var fields []Field
	for _, p := range props {
		// Skip non-maps with nested fields
		if !p.IsA("Map") && len(p.NestedProperties()) > 0 {
			continue
		}
		f := Field{
			Json:         p.IsJsonField(),
			ProviderOnly: p.ProviderOnly(),
		}
		lineage := p.Lineage()
		apiLineage := p.ApiLineage()
		// For maps (which all have nested children), instead of having an entry for the map field itself,
		// we want an entry that maps the `key` API field to the key_name of the map field.
		if p.IsA("Map") {
			lineage = append(lineage, p.KeyName)
			apiLineage = append(apiLineage, "key")
		}
		if !p.ProviderOnly() {
			f.ApiField = strings.Join(apiLineage, ".")
		}
		if p.ProviderOnly() || !IsDefaultLineage(lineage, apiLineage) {
			f.Field = strings.Join(lineage, ".")
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

// Returns true if the lineage is the default we'd expect for a field, and false otherwise.
// If any ancestor has a non-default lineage, this will return false.
func IsDefaultLineage(lineage, apiLineage []string) bool {
	if len(lineage) != len(apiLineage) {
		return false
	}
	for i, part := range lineage {
		apiPart := apiLineage[i]
		if part != google.Underscore(apiPart) {
			return false
		}
	}
	return true
}
