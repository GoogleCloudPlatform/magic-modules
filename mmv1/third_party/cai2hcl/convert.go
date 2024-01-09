package cai2hcl

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/caiasset"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/common"
	"go.uber.org/zap"
)

// Struct for options so that adding new options does not
// require updating function signatures all along the pipe.
type Options struct {
	ErrorLogger *zap.Logger
}

// Converts CAI Assets into HCL string.
func Convert(assets []*caiasset.Asset, options *Options) ([]byte, error) {
	if options == nil || options.ErrorLogger == nil {
		return nil, fmt.Errorf("logger is not initialized")
	}

	// Group resources from the same TF resource type for convert.
	// tf -> cai has 1:N mappings occasionally
	groups := make(map[string][]*caiasset.Asset)
	for _, asset := range assets {

		name, _ := AssetTypeToConverter[asset.Type]
		if name != "" {
			groups[name] = append(groups[name], asset)
		}
	}

	allBlocks := []*common.HCLResourceBlock{}
	for name, assets := range groups {
		converter, ok := ConverterMap[name]
		if !ok {
			continue
		}
		newBlocks, err := converter.Convert(assets)
		if err != nil {
			return nil, err
		}

		allBlocks = append(allBlocks, newBlocks...)
	}

	t, err := common.HclWriteBlocks(allBlocks)

	options.ErrorLogger.Debug(string(t))

	return t, err
}
