package toolkit

import (
	"sort"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tfplan2cai/converters"
)

// List of the supported Terraform resources
func ListSupportedTerraformResources() []string {
	resources := make([]string, 0, len(converters.ConverterMap))
	for k := range converters.ConverterMap {
		resources = append(resources, k)
	}
	sort.Strings(resources)
	return resources
}
