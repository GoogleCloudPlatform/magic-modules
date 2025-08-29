package models

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/caiasset"
)

// Converter interface for resources.
type Cai2hclConverter interface {
	// Convert turns asset into hcl blocks.
	Convert(asset caiasset.Asset) ([]*TerraformResourceBlock, error)
}
