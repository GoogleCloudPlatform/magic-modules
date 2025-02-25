package models

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/caiasset"
)

// Converter interface for resources.
type Converter interface {
	// Convert turns asset into hcl blocks.
	Convert(asset *caiasset.Asset) ([]*TerraformResourceBlock, error)
}
