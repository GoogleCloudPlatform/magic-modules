package networksecurity

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/cai2hcl/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/caiasset"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// BackendAuthenticationConfigAssetType is the CAI asset type name.
const BackendAuthenticationConfigAssetType string = "networksecurity.googleapis.com/BackendAuthenticationConfig"

// BackendAuthenticationConfigSchemaName is the TF resource schema name.
const BackendAuthenticationConfigSchemaName string = "google_network_security_backend_authentication_config"

// BackendAuthenticationConfigConverter for networksecurity backend authentication config resource.
type BackendAuthenticationConfigConverter struct {
	name   string
	schema map[string]*schema.Schema
}

// NewBackendAuthenticationConfigConverter returns an HCL converter.
func NewBackendAuthenticationConfigConverter(provider *schema.Provider) common.Converter {
	schema := provider.ResourcesMap[BackendAuthenticationConfigSchemaName].Schema

	return &BackendAuthenticationConfigConverter{
		name:   BackendAuthenticationConfigSchemaName,
		schema: schema,
	}
}

// Convert converts CAI assets to HCL resource blocks (Provider version: 6.45.0)
func (c *BackendAuthenticationConfigConverter) Convert(assets []*caiasset.Asset) ([]*common.HCLResourceBlock, error) {
	var blocks []*common.HCLResourceBlock
	var err error

	return blocks, err
}
