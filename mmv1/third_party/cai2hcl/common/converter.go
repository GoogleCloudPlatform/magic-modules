package common

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/caiasset"
	"github.com/zclconf/go-cty/cty"
)

// Converter interface for resources.
type Converter interface {
	// Convert turns assets into hcl blocks.
	Convert(asset []*caiasset.Asset) ([]*HCLResourceBlock, error)
}

// HCLResourceBlock identifies the HCL block's labels and content.
type HCLResourceBlock struct {
	Labels []string
	Value  cty.Value
}
