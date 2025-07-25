package networksecurity

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/cai2hcl/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/caiasset"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ServerTLSPolicyAssetType is the CAI asset type name.
const ServerTLSPolicyAssetType string = "networksecurity.googleapis.com/ServerTlsPolicy"

// ServerTLSPolicySchemaName is the TF resource schema name.
const ServerTLSPolicySchemaName string = "google_network_security_server_tls_policy"

// ServerTLSPolicyConverter for networksecurity server tls policy resource.
type ServerTLSPolicyConverter struct {
	name   string
	schema map[string]*schema.Schema
}

// NewServerTLSPolicyConverter returns an HCL converter.
func NewServerTLSPolicyConverter(provider *schema.Provider) common.Converter {
	schema := provider.ResourcesMap[ServerTLSPolicySchemaName].Schema

	return &ServerTLSPolicyConverter{
		name:   ServerTLSPolicySchemaName,
		schema: schema,
	}
}

// Convert converts CAI assets to HCL resource blocks.
func (c *ServerTLSPolicyConverter) Convert(assets []*caiasset.Asset) ([]*common.HCLResourceBlock, error) {
	panic("not implemented")
}
