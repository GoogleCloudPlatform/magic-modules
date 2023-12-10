package common

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tpg_provider "github.com/hashicorp/terraform-provider-google-beta/google-beta/provider"
)

// Function to create converter based on TF resource name and TF resource schema.
type ConverterFactory = func(name string, schema map[string]*schema.Schema) Converter

// Initializes converters from correpsonding converter factories.
func CreateConverters(converterFactories map[string]ConverterFactory) map[string]Converter {
	provider := tpg_provider.Provider()

	result := make(map[string]Converter, len(converterFactories))
	for name, factory := range converterFactories {
		result[name] = factory(name, provider.ResourcesMap[name].Schema)
	}

	return result
}
