package resourcemanager

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/cai2hcl/converters/utils"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/cai2hcl/models"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/caiasset"

	tfschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ProjectAssetType is the CAI asset type name for project.
const ProjectAssetType string = "cloudresourcemanager.googleapis.com/Project"

// ProjectSchemaName is the TF resource schema name for resourcemanager project.
const ProjectSchemaName string = "google_project"

// ProjectConverter for compute project resource.
type ProjectConverter struct {
	name   string
	schema map[string]*tfschema.Schema
}

// NewProjectConverter returns an HCL converter for compute project.
func NewProjectConverter(provider *tfschema.Provider) models.Converter {
	schema := provider.ResourcesMap[ProjectSchemaName].Schema

	return &ProjectConverter{
		name:   ProjectSchemaName,
		schema: schema,
	}
}

// Convert converts asset resource data.
func (c *ProjectConverter) Convert(asset caiasset.Asset) ([]*models.TerraformResourceBlock, error) {
	var blocks []*models.TerraformResourceBlock
	block, err := c.convertResourceData(asset)
	if err != nil {
		return nil, err
	}
	blocks = append(blocks, block)
	return blocks, nil
}

func (c *ProjectConverter) convertResourceData(asset caiasset.Asset) (*models.TerraformResourceBlock, error) {
	if asset.Resource == nil || asset.Resource.Data == nil {
		return nil, fmt.Errorf("asset resource data is nil")
	}

	assetResourceData := asset.Resource.Data

	hclData := make(map[string]interface{})
	hclData["name"] = assetResourceData["name"]
	hclData["project_id"] = assetResourceData["projectId"]
	hclData["labels"] = utils.RemoveTerraformAttributionLabel(assetResourceData["labels"])
	if strings.Contains(asset.Resource.Parent, "folders/") {
		hclData["folder_id"] = utils.ParseFieldValue(asset.Resource.Parent, "folders")
	} else if strings.Contains(asset.Resource.Parent, "organizations/") {
		hclData["org_id"] = utils.ParseFieldValue(asset.Resource.Parent, "organizations")
	}

	ctyVal, err := utils.MapToCtyValWithSchema(hclData, c.schema)
	if err != nil {
		return nil, err
	}
	return &models.TerraformResourceBlock{
		Labels: []string{c.name, assetResourceData["projectId"].(string)},
		Value:  ctyVal,
	}, nil
}
