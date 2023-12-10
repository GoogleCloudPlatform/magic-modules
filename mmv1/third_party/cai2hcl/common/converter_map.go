package common

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/caiasset"
)

// Configuration object to map CAI Assets to Converters.
type ConverterMap struct {
	// Mapping from CAI Asset Type to converter name.
	AssetTypeToConverter map[string]string

	// Mapping from converter name and converter.
	Converters map[string]Converter
}

// ConvertMap implements Converter interface to be able to convert any Asset type.
func (c ConverterMap) Convert(assets []*caiasset.Asset) ([]*HCLResourceBlock, error) {
	// Group resources from the same tf resource type for convert.
	// tf -> cai has 1:N mappings occasionally
	groups := make(map[string][]*caiasset.Asset)
	for _, asset := range assets {

		name, _ := c.AssetTypeToConverter[asset.Type]
		if name != "" {
			groups[name] = append(groups[name], asset)
		}
	}

	allBlocks := []*HCLResourceBlock{}
	for name, assets := range groups {
		converter, ok := c.Converters[name]
		if !ok {
			continue
		}
		newBlocks, err := converter.Convert(assets)
		if err != nil {
			return nil, err
		}

		allBlocks = append(allBlocks, newBlocks...)
	}

	return allBlocks, nil
}
