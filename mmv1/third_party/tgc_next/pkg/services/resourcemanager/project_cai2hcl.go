package resourcemanager

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/cai2hcl/converters/utils"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/cai2hcl/models"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tgcresource"

	tfschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ProjectCai2hclConverter for compute project resource.
type ProjectCai2hclConverter struct {
	name   string
	schema map[string]*tfschema.Schema
}

// NewProjectConverter returns an HCL converter for compute project.
func NewProjectCai2hclConverter(provider *tfschema.Provider) models.Cai2hclConverter {
	schema := provider.ResourcesMap[ProjectSchemaName].Schema

	return &ProjectCai2hclConverter{
		name:   ProjectSchemaName,
		schema: schema,
	}
}

// Convert converts asset resource data.
<<<<<<<< HEAD:mmv1/third_party/tgc_next/pkg/cai2hcl/converters/services/resourcemanager/project.go
func (c *ProjectConverter) Convert(asset caiasset.Asset) ([]*models.TerraformResourceBlock, error) {
========
func (c *ProjectCai2hclConverter) Convert(asset caiasset.Asset) ([]*models.TerraformResourceBlock, error) {
>>>>>>>> 4dd5624b9bd5d9fb7d39acf94deb53127832f1d1:mmv1/third_party/tgc_next/pkg/services/resourcemanager/project_cai2hcl.go
	var blocks []*models.TerraformResourceBlock
	block, err := c.convertResourceData(asset)
	if err != nil {
		return nil, err
	}
	blocks = append(blocks, block)
	return blocks, nil
}

<<<<<<<< HEAD:mmv1/third_party/tgc_next/pkg/cai2hcl/converters/services/resourcemanager/project.go
func (c *ProjectConverter) convertResourceData(asset caiasset.Asset) (*models.TerraformResourceBlock, error) {
========
func (c *ProjectCai2hclConverter) convertResourceData(asset caiasset.Asset) (*models.TerraformResourceBlock, error) {
>>>>>>>> 4dd5624b9bd5d9fb7d39acf94deb53127832f1d1:mmv1/third_party/tgc_next/pkg/services/resourcemanager/project_cai2hcl.go
	if asset.Resource == nil || asset.Resource.Data == nil {
		return nil, fmt.Errorf("asset resource data is nil")
	}

	assetResourceData := asset.Resource.Data

	hclData := make(map[string]interface{})
	hclData["name"] = assetResourceData["name"]
	hclData["project_id"] = assetResourceData["projectId"]
	hclData["labels"] = tgcresource.RemoveTerraformAttributionLabel(assetResourceData["labels"])
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
