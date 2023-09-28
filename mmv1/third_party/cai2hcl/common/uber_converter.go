package common

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/cai2hcl/common/matchers"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/caiasset"
)

// Converter which aggregates all service-specific converters in the same interface.
type UberConverter struct {
	// Mapping between asset type (i.e. compute.googleapis.com/Instance) to collection of matchers.
	// Collection of asset name formats (i.e. projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/forwardingRules)
	// together with corresponding converter name.
	ConverterByAssetType map[string]string

	ConverterMatchersByAssetType map[string][]matchers.ConverterMatcher

	// Mapping between converter name and converter constructor.
	Converters map[string]Converter
}

// Convert assets of any of known types to the list of HCL blocks.
func (c UberConverter) Convert(assets []*caiasset.Asset) ([]*HCLResourceBlock, error) {
	// Group resources from the same tf resource type for convert.
	// tf -> cai has 1:N mappings occasionally
	groups := make(map[string][]*caiasset.Asset)
	for _, asset := range assets {

		name, found := c.ConverterByAssetType[asset.Type]
		if !found {
			matchers, _ := c.ConverterMatchersByAssetType[asset.Type]

			for _, matcher := range matchers {
				if matcher.Match(asset) {
					name = matcher.GetConverterName()
				}
			}
		}

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
