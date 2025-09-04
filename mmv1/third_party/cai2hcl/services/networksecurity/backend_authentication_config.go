package networksecurity

import (
	"errors"
	"fmt"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/cai2hcl/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/caiasset"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	netsecapi "google.golang.org/api/networksecurity/v1"
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

// Convert converts CAI assets to HCL resource blocks (Provider version: 7.0.1)
func (c *BackendAuthenticationConfigConverter) Convert(assets []*caiasset.Asset) ([]*common.HCLResourceBlock, error) {
	var blocks []*common.HCLResourceBlock
	var err error

	for _, asset := range assets {
		if asset == nil {
			continue
		} else if asset.Resource == nil || asset.Resource.Data == nil {
			return nil, fmt.Errorf("INVALID_ARGUMENT: Asset resource data is nil")
		} else if asset.Type != BackendAuthenticationConfigAssetType {
			return nil, fmt.Errorf("INVALID_ARGUMENT: Expected asset of type %s, but received %s", BackendAuthenticationConfigAssetType, asset.Type)
		}
		block, errConvert := c.convertResourceData(asset)
		blocks = append(blocks, block)
		if errConvert != nil {
			err = errors.Join(err, errConvert)
		}
	}
	return blocks, err
}

func (c *BackendAuthenticationConfigConverter) convertResourceData(asset *caiasset.Asset) (*common.HCLResourceBlock, error) {
	if asset == nil || asset.Resource == nil || asset.Resource.Data == nil {
		return nil, fmt.Errorf("INVALID_ARGUMENT: Asset resource data is nil")
	}

	hcl, _ := flattenBackendAuthenticationConfig(asset.Resource)

	ctyVal, err := common.MapToCtyValWithSchema(hcl, c.schema)
	if err != nil {
		return nil, err
	}

	resourceName := hcl["name"].(string)
	return &common.HCLResourceBlock{
		Labels: []string{c.name, resourceName},
		Value:  ctyVal,
	}, nil
}

func flattenBackendAuthenticationConfig(resource *caiasset.AssetResource) (map[string]any, error) {
	result := make(map[string]any)

	var backendAuthenticationConfig *netsecapi.BackendAuthenticationConfig
	if err := common.DecodeJSON(resource.Data, &backendAuthenticationConfig); err != nil {
		return nil, err
	}

	result["name"] = flattenName(backendAuthenticationConfig.Name)
	result["labels"] = backendAuthenticationConfig.Labels
	result["description"] = backendAuthenticationConfig.Description
	result["client_certificate"] = backendAuthenticationConfig.ClientCertificate
	result["trust_config"] = backendAuthenticationConfig.TrustConfig
	result["well_known_roots"] = backendAuthenticationConfig.WellKnownRoots
	result["project"] = flattenProjectName(backendAuthenticationConfig.Name)

	result["location"] = resource.Location

	return result, nil
}
