package certificatemanager

import (
	"errors"
	"fmt"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/cai2hcl/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/caiasset"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	certificatemanagerapi "google.golang.org/api/certificatemanager/v1"
	"strings"
)

// CertificateAssetType is the CAI asset type name.
const CertificateAssetType string = "certificatemanager.googleapis.com/Certificate"

// CertificateSchemaName is the TF resource schema name.
const CertificateSchemaName string = "google_certificate_manager_certificate"

// CertificateConverter for certificatemanager Certificate resource
type CertificateConverter struct {
	name   string
	schema map[string]*schema.Schema
}

// NewCertificateConverter returns an HCL converter
func NewCertificateConverter(provider *schema.Provider) common.Converter {
	schema := provider.ResourcesMap[CertificateSchemaName].Schema

	return &CertificateConverter{
		name:   CertificateSchemaName,
		schema: schema,
	}
}

// Convert converts CAI assets to HCL resource blocks (Provider version: 6.47.0)
func (c *CertificateConverter) Convert(assets []*caiasset.Asset) ([]*common.HCLResourceBlock, error) {
	var blocks []*common.HCLResourceBlock
	var err error

	for _, asset := range assets {
		if asset == nil {
			continue
		} else if asset.Resource == nil || asset.Resource.Data == nil {
			return nil, fmt.Errorf("INVALID_ARGUMENT: Asset resource data is nil")
		} else if asset.Type != CertificateAssetType {
			return nil, fmt.Errorf("INVALID_ARGUMENT: Expected asset of type %s, but received %s", CertificateAssetType, asset.Type)
		}
		block, errConvert := c.convertResourceData(asset)
		blocks = append(blocks, block)
		if errConvert != nil {
			err = errors.Join(err, errConvert)
		}
	}
	return blocks, err
}

func (c *CertificateConverter) convertResourceData(asset *caiasset.Asset) (*common.HCLResourceBlock, error) {
	if asset == nil || asset.Resource == nil || asset.Resource.Data == nil {
		return nil, fmt.Errorf("INVALID_ARGUMENT: Asset resource data is nil")
	}

	hcl, _ := flattenCertificate(asset.Resource)

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

func flattenCertificate(resource *caiasset.AssetResource) (map[string]any, error) {
	result := make(map[string]any)

	var certificate *certificatemanagerapi.Certificate
	if err := common.DecodeJSON(resource.Data, &certificate); err != nil {
		return nil, err
	}

	result["name"] = flattenName(certificate.Name)
	result["description"] = certificate.Description
	result["labels"] = certificate.Labels
	result["scope"] = certificate.Scope
	result["self_managed"] = flattenSelfManaged(certificate.SelfManaged, certificate.PemCertificate)
	result["managed"] = flattenManaged(certificate.Managed)
	result["project"] = flattenProjectName(certificate.Name)

	result["location"] = resource.Location

	return result, nil
}

func flattenName(name string) string {
	tokens := strings.Split(name, "/")
	return tokens[len(tokens)-1]
}

func flattenSelfManaged(selfManaged *certificatemanagerapi.SelfManagedCertificate, pemCertificate string) []map[string]any {
	if selfManaged == nil {
		return nil
	}

	result := make(map[string]any)
	result["pem_certificate"] = pemCertificate
	result["pem_private_key"] = "<private_key>"

	return []map[string]any{result}
}

func flattenManaged(managed *certificatemanagerapi.ManagedCertificate) []map[string]any {
	if managed == nil {
		return nil
	}

	result := make(map[string]any)
	result["domains"] = managed.Domains
	result["dns_authorizations"] = managed.DnsAuthorizations
	result["issuance_config"] = managed.IssuanceConfig

	return []map[string]any{result}
}

func flattenProjectName(name string) string {
	tokens := strings.Split(name, "/")
	if len(tokens) < 2 || tokens[0] != "projects" {
		return ""
	}
	return tokens[1]
}
