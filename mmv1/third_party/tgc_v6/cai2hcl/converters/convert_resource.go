package converters

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/cai2hcl/models"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/caiasset"
)

func ConvertResource(asset *caiasset.Asset) ([]*models.TerraformResourceBlock, error) {
	converter, ok := ConverterMap[asset.Type]
	if !ok {
		return nil, nil
	}
	return converter.Convert(asset)
}
