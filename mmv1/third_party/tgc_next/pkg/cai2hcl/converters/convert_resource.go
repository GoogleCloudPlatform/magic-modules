package converters

import (
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/cai2hcl/models"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/caiasset"
)

func ConvertResource(asset caiasset.Asset) ([]*models.TerraformResourceBlock, error) {
	converters, ok := ConverterMap[asset.Type]
	if !ok || len(converters) == 0 {
		return nil, nil
	}

	var converter models.Cai2hclConverter
	// Normally, one asset type has only one converter.
	if len(converters) == 1 {
		for _, converter = range converters {
			return converter.Convert(asset)
		}
	}

	// Edge cases
	if asset.Type == "compute.googleapis.com/Autoscaler" {
		if strings.Contains(asset.Name, "/zones/") {
			converter = ConverterMap[asset.Type]["ComputeAutoscaler"]
		} else {
			converter = ConverterMap[asset.Type]["ComputeRegionAutoscaler"]
		}
	}
	return converter.Convert(asset)
}
