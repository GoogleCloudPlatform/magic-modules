package resourcemanager

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/caiasset"

	tfschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zclconf/go-cty/cty"
	"google.golang.org/api/cloudresourcemanager/v1"
)

// ProjectAssetType is the CAI asset type name for project.
const ProjectAssetType string = "cloudresourcemanager.googleapis.com/Project"

// ProjectAssetType is the CAI asset type name for project.
const ProjectBillingAssetType string = "cloudbilling.googleapis.com/ProjectBillingInfo"

// ProjectSchemaName is the TF resource schema name for resourcemanager project.
const ProjectSchemaName string = "google_project"

// ProjectConverter for compute project resource.
type ProjectConverter struct {
	name     string
	schema   map[string]*tfschema.Schema
	billings map[string]string
}

// NewProjectConverter returns an HCL converter for compute project.
func NewProjectConverter(provider *tfschema.Provider) common.Converter {
	schema := provider.ResourcesMap[ProjectSchemaName].Schema

	return &ProjectConverter{
		name:     ProjectSchemaName,
		schema:   schema,
		billings: make(map[string]string),
	}
}

// Convert converts asset resource data.
func (c *ProjectConverter) Convert(assets []*caiasset.Asset) ([]*common.HCLResourceBlock, error) {
	// process billing info
	for _, asset := range assets {
		if asset == nil {
			continue
		}
		if asset.Type == "cloudbilling.googleapis.com/ProjectBillingInfo" {
			project := common.ParseFieldValue(asset.Name, "projects")
			projectAssetName := fmt.Sprintf("//cloudresourcemanager.googleapis.com/projects/%s", project)
			c.billings[projectAssetName] = c.convertBilling(asset)
		}
	}

	var blocks []*common.HCLResourceBlock
	for _, asset := range assets {
		if asset == nil {
			continue
		}
		if asset.Type == "cloudbilling.googleapis.com/ProjectBillingInfo" {
			continue
		}
		if asset.IAMPolicy != nil {
			iamBlock, err := c.convertIAM(asset)
			if err != nil {
				return nil, err
			}
			blocks = append(blocks, iamBlock)
		}
		if asset.Resource != nil && asset.Resource.Data != nil {
			block, err := c.convertResourceData(asset)
			if err != nil {
				return nil, err
			}
			blocks = append(blocks, block)
		}
	}
	return blocks, nil
}

func (c *ProjectConverter) convertIAM(asset *caiasset.Asset) (*common.HCLResourceBlock, error) {
	if asset.IAMPolicy == nil {
		return nil, fmt.Errorf("asset IAM policy is nil")
	}

	project := common.ParseFieldValue(asset.Name, "projects")
	policyData, err := json.Marshal(asset.IAMPolicy)
	if err != nil {
		return nil, err
	}

	return &common.HCLResourceBlock{
		Labels: []string{
			c.name + "_iam_policy",
			project + "_iam_policy",
		},
		Value: cty.ObjectVal(map[string]cty.Value{
			"project":     cty.StringVal(project),
			"policy_data": cty.StringVal(string(policyData)),
		}),
	}, nil
}

func (c *ProjectConverter) convertBilling(asset *caiasset.Asset) string {
	if asset != nil && asset.Resource != nil && asset.Resource.Data != nil {
		return strings.TrimPrefix(asset.Resource.Data["billingAccountName"].(string), "billingAccounts/")
	}
	return ""
}

func (c *ProjectConverter) convertResourceData(asset *caiasset.Asset) (*common.HCLResourceBlock, error) {
	if asset == nil || asset.Resource == nil || asset.Resource.Data == nil {
		return nil, fmt.Errorf("asset resource data is nil")
	}
	var project *cloudresourcemanager.Project
	if err := common.DecodeJSON(asset.Resource.Data, &project); err != nil {
		return nil, err
	}

	hclData := make(map[string]interface{})
	hclData["name"] = project.Name
	hclData["project_id"] = project.ProjectId
	hclData["labels"] = project.Labels
	if strings.Contains(asset.Resource.Parent, "folders/") {
		hclData["folder_id"] = common.ParseFieldValue(asset.Resource.Parent, "folders")
	} else if strings.Contains(asset.Resource.Parent, "organizations/") {
		hclData["org_id"] = common.ParseFieldValue(asset.Resource.Parent, "organizations")
	}

	if billingAccount, ok := c.billings[asset.Name]; ok {
		hclData["billing_account"] = billingAccount
	}

	ctyVal, err := common.MapToCtyValWithSchema(hclData, c.schema)
	if err != nil {
		return nil, err
	}
	return &common.HCLResourceBlock{
		Labels: []string{c.name, project.ProjectId},
		Value:  ctyVal,
	}, nil
}
