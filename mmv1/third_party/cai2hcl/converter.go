package cai2hcl

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/caiasset"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/cai2hcl/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/cai2hcl/services"
	"go.uber.org/zap"
)

// Struct for options so that adding new options does not
// require updating function signatures all along the pipe.
type Options struct {
	ErrorLogger *zap.Logger
}

// Converts CAI Assets into HCL.
func Convert(assets []*caiasset.Asset, options *Options) ([]byte, error) {
	if options == nil || options.ErrorLogger == nil {
		return nil, fmt.Errorf("logger is not initialized")
	}

	servicesConverter := services.UberConverter

	blocks, err := servicesConverter.Convert(assets)
	if err != nil {
		return nil, err
	}

	t, err := common.HclWriteBlocks(blocks)

	options.ErrorLogger.Debug(string(t))

	return t, err
}
