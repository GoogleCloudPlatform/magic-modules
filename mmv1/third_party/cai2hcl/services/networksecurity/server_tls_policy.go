package networksecurity

import (
	"errors"
	"fmt"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/cai2hcl/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/caiasset"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	netsecapi "google.golang.org/api/networksecurity/v1"
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

// Convert converts CAI assets to HCL resource blocks (Provider version: 6.45.0)
func (c *ServerTLSPolicyConverter) Convert(assets []*caiasset.Asset) ([]*common.HCLResourceBlock, error) {
	var blocks []*common.HCLResourceBlock
	var err error

	for _, asset := range assets {
		if asset == nil {
			continue
		} else if asset.Resource == nil || asset.Resource.Data == nil {
			return nil, fmt.Errorf("INVALID_ARGUMENT: Asset resource data is nil")
		} else if asset.Type != ServerTLSPolicyAssetType {
			return nil, fmt.Errorf("INVALID_ARGUMENT: Expected asset of type %s, but received %s", ServerTLSPolicyAssetType, asset.Type)
		}
		block, errConvert := c.convertResourceData(asset)
		blocks = append(blocks, block)
		if errConvert != nil {
			err = errors.Join(err, errConvert)
		}
	}
	return blocks, err
}

func (c *ServerTLSPolicyConverter) convertResourceData(asset *caiasset.Asset) (*common.HCLResourceBlock, error) {
	if asset == nil || asset.Resource == nil || asset.Resource.Data == nil {
		return nil, fmt.Errorf("INVALID_ARGUMENT: Asset resource data is nil")
	}

	hcl, _ := flattenServerTLSPolicy(asset.Resource)

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

func flattenServerTLSPolicy(resource *caiasset.AssetResource) (map[string]any, error) {
	result := make(map[string]any)

	var serverTLSPolicy *netsecapi.ServerTlsPolicy
	if err := common.DecodeJSON(resource.Data, &serverTLSPolicy); err != nil {
		return nil, err
	}

	result["name"] = flattenName(serverTLSPolicy.Name)
	result["labels"] = serverTLSPolicy.Labels
	result["description"] = serverTLSPolicy.Description
	result["allow_open"] = serverTLSPolicy.AllowOpen
	result["server_certificate"] = flattenServerCertificate(serverTLSPolicy.ServerCertificate)
	result["mtls_policy"] = flattenMTLSPolicy(serverTLSPolicy.MtlsPolicy)
	result["project"] = flattenProjectName(serverTLSPolicy.Name)

	result["location"] = resource.Location

	return result, nil
}

func flattenServerCertificate(certificate *netsecapi.GoogleCloudNetworksecurityV1CertificateProvider) []map[string]any {
	if certificate == nil {
		return nil
	}

	result := make(map[string]any)
	result["certificate_provider_instance"] = flattenCertificateProviderInstance(certificate.CertificateProviderInstance)
	result["grpc_endpoint"] = flattenGrpcEndpoint(certificate.GrpcEndpoint)

	return []map[string]any{result}
}

func flattenMTLSPolicy(policy *netsecapi.MTLSPolicy) []map[string]any {
	if policy == nil {
		return nil
	}

	result := make(map[string]any)
	result["client_validation_mode"] = policy.ClientValidationMode
	result["client_validation_trust_config"] = policy.ClientValidationTrustConfig
	result["client_validation_ca"] = flattenClientValidationCA(policy.ClientValidationCa)

	return []map[string]any{result}
}

func flattenCertificateProviderInstance(instance *netsecapi.CertificateProviderInstance) []map[string]any {
	if instance == nil {
		return nil
	}

	result := make(map[string]any)
	result["plugin_instance"] = instance.PluginInstance

	return []map[string]any{result}
}

func flattenGrpcEndpoint(endpoint *netsecapi.GoogleCloudNetworksecurityV1GrpcEndpoint) []map[string]any {
	if endpoint == nil {
		return nil
	}

	result := make(map[string]any)
	result["target_uri"] = endpoint.TargetUri

	return []map[string]any{result}
}

func flattenClientValidationCA(cas []*netsecapi.ValidationCA) []map[string]any {
	if cas == nil {
		return nil
	}

	result := make([]map[string]any, 0, len(cas))

	for _, ca := range cas {
		converted := map[string]any{
			"certificate_provider_instance": flattenCertificateProviderInstance(ca.CertificateProviderInstance),
			"grpc_endpoint":                 flattenGrpcEndpoint(ca.GrpcEndpoint),
		}
		result = append(result, converted)
	}

	return result
}
