package common

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/provider"
)

// Function initializing a converter from TF resource name and TF resource schema.
type ConverterFactory = func(name string, schema map[string]*schema.Schema) Converter

// Initializes map of converters.
func CreateConverterMap(converterFactories map[string]ConverterFactory) map[string]Converter {
	tpgProvider := tpg.Provider()

	result := make(map[string]Converter, len(converterFactories))
	for name, factory := range converterFactories {
		result[name] = factory(name, tpgProvider.ResourcesMap[name].Schema)
	}

	return result
}
