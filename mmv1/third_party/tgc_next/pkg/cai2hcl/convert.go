package cai2hcl

import (
	"bytes"
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/cai2hcl/converters"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/cai2hcl/models"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/caiasset"
	"go.uber.org/zap"
)

// Struct for options so that adding new options does not
// require updating function signatures all along the pipe.
type Options struct {
	ErrorLogger *zap.Logger
}

// Converts CAI Assets into HCL string.
func Convert(assets []caiasset.Asset, options *Options) ([]byte, error) {
	if options == nil || options.ErrorLogger == nil {
		return nil, fmt.Errorf("logger is not initialized")
	}

	// TODO: add resolvers to resolve the assets into single resource assets

	allResourceBytes, err := converters.ConvertResource(assets, &models.Options{})
	if err != nil {
		return nil, err
	}

	return bytes.Join(allResourceBytes, []byte("\n")), nil
}
