package models

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/caiasset"
)

// Converter interface for resources.
type Cai2hclConverter interface {
	// Convert turns asset into hcl blocks.
	Convert(assets []caiasset.Asset, options *Options) ([]*TerraformResourceBlock, error)
}
