package common

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/caiasset"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zclconf/go-cty/cty"
)

// Converter interface for resources.
type Converter interface {
	// Convert turns assets into hcl blocks.
	Convert(asset []*caiasset.Asset) ([]*HCLResourceBlock, error)
}

// Function initializing a converter from TF resource name and TF resource schema.
type ConverterFactory = func(name string, schema map[string]*schema.Schema) Converter

// HCLResourceBlock identifies the HCL block's labels and content.
type HCLResourceBlock struct {
	Labels []string
	Value  cty.Value
}
