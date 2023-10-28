package common

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tpg_provider "github.com/hashicorp/terraform-provider-google-beta/google-beta/provider"
)

// Function initializing a converter from TF resource name and TF resource schema.
type ConverterFactory = func(name string, schema map[string]*schema.Schema) Converter

// Initializes map of converters.
func CreateConverterMap(converterFactories map[string]ConverterFactory) map[string]Converter {
	provider := tpg_provider.Provider()

	result := make(map[string]Converter, len(converterFactories))
	for name, factory := range converterFactories {
		result[name] = factory(name, provider.ResourcesMap[name].Schema)
	}

	return result
}
