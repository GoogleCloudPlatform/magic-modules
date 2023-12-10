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

	converter := ConverterMap

	blocks, err := converter.Convert(assets)
	if err != nil {
		return nil, err
	}

	t, err := common.HclWriteBlocks(blocks)

	options.ErrorLogger.Debug(string(t))

	return t, err
}
