package compute

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/cai2hcl/converters/utils"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/cai2hcl/models"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tgcresource"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tpgresource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ComputeResizeRequestCai2hclConverter struct {
	name   string
	schema map[string]*schema.Schema
}

func NewComputeResizeRequestCai2hclConverter(provider *schema.Provider) models.Cai2hclConverter {
	schema := provider.ResourcesMap[ComputeResizeRequestSchemaName].Schema

	return &ComputeResizeRequestCai2hclConverter{
		name:   ComputeResizeRequestSchemaName,
		schema: schema,
	}
}

func (c *ComputeResizeRequestCai2hclConverter) Convert(assets []caiasset.Asset, options *models.ResourceConverterOptions) ([]*models.TerraformResourceBlock, error) {
	if len(assets) == 0 {
		return nil, nil
	}

	var blocks []*models.TerraformResourceBlock
	for _, asset := range assets {
		if asset.Resource == nil || asset.Resource.Data == nil {
			continue
		}
		data := asset.Resource.Data
		name, _ := data["name"].(string)
		if name == "" {
			name = tpgresource.GetResourceNameFromSelfLink(asset.Name)
		}

		hclData := make(map[string]any)
		if name != "" {
			hclData["name"] = name
		}
		if igm, ok := data["instanceGroupManager"].(string); ok && igm != "" {
			hclData["instance_group_manager"] = tpgresource.GetResourceNameFromSelfLink(igm)
		} else if igmName := tgcresource.ParseFieldValue(asset.Name, "instanceGroupManagers"); igmName != "" {
			hclData["instance_group_manager"] = igmName
		}
		if resizeBy, ok := data["resizeBy"].(float64); ok {
			hclData["resize_by"] = int(resizeBy)
		} else if count, ok := data["count"].(float64); ok {
			hclData["resize_by"] = int(count)
		}
		if description, ok := data["description"].(string); ok && description != "" {
			hclData["description"] = description
		}
		if zone, ok := data["zone"].(string); ok && zone != "" {
			hclData["zone"] = zone
		} else if zone := tgcresource.ParseFieldValue(asset.Name, "zones"); zone != "" {
			hclData["zone"] = zone
		}
		if reqDuration := flattenResizeRequestRequestedRunDuration(data["requestedRunDuration"]); len(reqDuration) > 0 {
			hclData["requested_run_duration"] = reqDuration
		}

		ctyVal, err := utils.MapToCtyValWithSchema(hclData, c.schema)
		if err != nil {
			return nil, fmt.Errorf("error converting map to cty.Value: %w", err)
		}

		block := &models.TerraformResourceBlock{
			Labels: []string{c.name, name},
			Value:  ctyVal,
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}

func flattenResizeRequestRequestedRunDuration(v any) []map[string]any {
	if v == nil {
		return nil
	}
	m, ok := v.(map[string]any)
	if !ok || len(m) == 0 {
		return nil
	}
	item := make(map[string]any)
	if sec, ok := m["seconds"].(string); ok && sec != "" {
		item["seconds"] = sec
	} else if secNum, ok := m["seconds"].(float64); ok {
		item["seconds"] = fmt.Sprintf("%.0f", secNum)
	}
	if nanos, ok := m["nanos"].(float64); ok {
		item["nanos"] = int(nanos)
	} else if nanosStr, ok := m["nanos"].(string); ok && nanosStr != "" {
		var n int
		fmt.Sscanf(nanosStr, "%d", &n)
		item["nanos"] = n
	}
	if len(item) == 0 {
		return nil
	}
	return []map[string]any{item}
}
