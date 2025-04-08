package cai2hcl

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/cai2hcl/converters"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/cai2hcl/models"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/caiasset"
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

	// TODO: add resolvers to resolve the assets into single resource assets

	allBlocks := []*models.TerraformResourceBlock{}
	for _, asset := range assets {
		newBlocks, err := converters.ConvertResource(asset)
		if err != nil {
			return nil, err
		}

		if newBlocks != nil {
			allBlocks = append(allBlocks, newBlocks...)
		}
	}

	t, err := models.HclWriteBlocks(allBlocks)

	options.ErrorLogger.Debug(string(t))

	return t, err
}
