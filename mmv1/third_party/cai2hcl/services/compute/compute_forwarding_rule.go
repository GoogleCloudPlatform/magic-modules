package compute

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/caiasset"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

func forwardingRuleCustomizeDiff(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
	log.Println("[DEBUG] [PSC] Reached forwardingRuleCustomizeDiff function")

	// if target is not a string it's not set so no PSC connection
	if target, ok := diff.Get("target").(string); ok {
		if strings.Contains(target, "/serviceAttachments/") {
			recreateClosedPsc, _ := diff.Get("recreate_closed_psc").(bool)
			if pscConnectionStatus, ok := diff.Get("psc_connection_status").(string); ok && recreateClosedPsc && pscConnectionStatus == "CLOSED" {
				// https://discuss.hashicorp.com/t/force-new-resource-based-on-api-read-difference/29759/6
				diff.SetNewComputed("psc_connection_status")
				diff.ForceNew("psc_connection_status")
			}
		}
	}
	return nil
}

// ComputeForwardingRuleAssetType is the CAI asset type name.
const ComputeForwardingRuleAssetType string = "compute.googleapis.com/ForwardingRule"

// ComputeForwardingRuleSchemaName is a TF resource schema name.
const ComputeForwardingRuleSchemaName string = "google_compute_forwarding_rule"

type ComputeForwardingRuleConverter struct {
	name   string
	schema map[string]*schema.Schema
}

// NewComputeForwardingRuleConverter returns an HCL converter for compute instance.
func NewComputeForwardingRuleConverter(provider *schema.Provider) common.Converter {
	schema := provider.ResourcesMap[ComputeForwardingRuleSchemaName].Schema

	return &ComputeForwardingRuleConverter{
		name:   ComputeForwardingRuleSchemaName,
		schema: schema,
	}
}

func (c *ComputeForwardingRuleConverter) Convert(assets []*caiasset.Asset) ([]*common.HCLResourceBlock, error) {
	var blocks []*common.HCLResourceBlock
	config := common.NewConfig()

	for _, asset := range assets {
		if asset == nil {
			continue
		}
		if asset.Resource != nil && asset.Resource.Data != nil {
			block, err := c.convertResourceData(asset, config)
			if err != nil {
				return nil, err
			}
			blocks = append(blocks, block)
		}
	}
	return blocks, nil
}

func (c *ComputeForwardingRuleConverter) convertResourceData(asset *caiasset.Asset, config *transport_tpg.Config) (*common.HCLResourceBlock, error) {
	if asset == nil || asset.Resource == nil || asset.Resource.Data == nil {
		return nil, fmt.Errorf("asset resource data is nil")
	}

	assetResourceData := asset.Resource.Data

	hcl, _ := flattenComputeForwardingRule(assetResourceData, config)

	ctyVal, err := common.MapToCtyValWithSchema(hcl, c.schema)
	if err != nil {
		return nil, err
	}

	resourceName := assetResourceData["name"].(string)

	return &common.HCLResourceBlock{
		Labels: []string{c.name, resourceName},
		Value:  ctyVal,
	}, nil
}

func flattenComputeForwardingRule(resource map[string]interface{}, config *transport_tpg.Config) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	var resource_data *schema.ResourceData = nil

	result["creation_timestamp"] = flattenComputeForwardingRuleCreationTimestamp(resource["creationTimestamp"], resource_data, config)
	result["is_mirroring_collector"] = flattenComputeForwardingRuleIsMirroringCollector(resource["isMirroringCollector"], resource_data, config)
	result["psc_connection_id"] = flattenComputeForwardingRulePscConnectionId(resource["pscConnectionId"], resource_data, config)
	result["psc_connection_status"] = flattenComputeForwardingRulePscConnectionStatus(resource["pscConnectionStatus"], resource_data, config)
	result["description"] = flattenComputeForwardingRuleDescription(resource["description"], resource_data, config)
	result["ip_address"] = flattenComputeForwardingRuleIPAddress(resource["IPAddress"], resource_data, config)
	result["ip_protocol"] = flattenComputeForwardingRuleIPProtocol(resource["IPProtocol"], resource_data, config)
	result["backend_service"] = flattenComputeForwardingRuleBackendService(resource["backendService"], resource_data, config)
	result["load_balancing_scheme"] = flattenComputeForwardingRuleLoadBalancingScheme(resource["loadBalancingScheme"], resource_data, config)
	result["name"] = flattenComputeForwardingRuleName(resource["name"], resource_data, config)
	result["network"] = flattenComputeForwardingRuleNetwork(resource["network"], resource_data, config)
	result["port_range"] = flattenComputeForwardingRulePortRange(resource["portRange"], resource_data, config)
	result["ports"] = flattenComputeForwardingRulePorts(resource["ports"], resource_data, config)
	result["subnetwork"] = flattenComputeForwardingRuleSubnetwork(resource["subnetwork"], resource_data, config)
	result["target"] = flattenComputeForwardingRuleTarget(resource["target"], resource_data, config)
	result["allow_global_access"] = flattenComputeForwardingRuleAllowGlobalAccess(resource["allowGlobalAccess"], resource_data, config)
	result["labels"] = flattenComputeForwardingRuleLabels(resource["labels"], resource_data, config)
	result["label_fingerprint"] = flattenComputeForwardingRuleLabelFingerprint(resource["labelFingerprint"], resource_data, config)
	result["all_ports"] = flattenComputeForwardingRuleAllPorts(resource["allPorts"], resource_data, config)
	result["network_tier"] = flattenComputeForwardingRuleNetworkTier(resource["networkTier"], resource_data, config)
	result["service_directory_registrations"] = flattenComputeForwardingRuleServiceDirectoryRegistrations(resource["serviceDirectoryRegistrations"], resource_data, config)
	result["service_label"] = flattenComputeForwardingRuleServiceLabel(resource["serviceLabel"], resource_data, config)
	result["service_name"] = flattenComputeForwardingRuleServiceName(resource["serviceName"], resource_data, config)
	result["source_ip_ranges"] = flattenComputeForwardingRuleSourceIpRanges(resource["sourceIpRanges"], resource_data, config)
	result["base_forwarding_rule"] = flattenComputeForwardingRuleBaseForwardingRule(resource["baseForwardingRule"], resource_data, config)
	result["allow_psc_global_access"] = flattenComputeForwardingRuleAllowPscGlobalAccess(resource["allowPscGlobalAccess"], resource_data, config)
	result["ip_version"] = flattenComputeForwardingRuleIpVersion(resource["ipVersion"], resource_data, config)
	result["terraform_labels"] = flattenComputeForwardingRuleTerraformLabels(resource["labels"], resource_data, config)
	result["effective_labels"] = flattenComputeForwardingRuleEffectiveLabels(resource["labels"], resource_data, config)
	result["region"] = flattenComputeForwardingRuleRegion(resource["region"], resource_data, config)

	return result, nil
}

func flattenComputeForwardingRuleCreationTimestamp(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRuleIsMirroringCollector(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRulePscConnectionId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRulePscConnectionStatus(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRuleDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRuleIPAddress(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRuleIPProtocol(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRuleBackendService(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.ConvertSelfLinkToV1(v.(string))
}

func flattenComputeForwardingRuleLoadBalancingScheme(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRuleName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRuleNetwork(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.ConvertSelfLinkToV1(v.(string))
}

func flattenComputeForwardingRulePortRange(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRulePorts(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return schema.NewSet(schema.HashString, v.([]interface{}))
}

func flattenComputeForwardingRuleSubnetwork(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.ConvertSelfLinkToV1(v.(string))
}

func flattenComputeForwardingRuleTarget(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRuleAllowGlobalAccess(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRuleLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}

	transformed := make(map[string]interface{})
	if l, ok := d.GetOkExists("labels"); ok {
		for k := range l.(map[string]interface{}) {
			transformed[k] = v.(map[string]interface{})[k]
		}
	}

	return transformed
}

func flattenComputeForwardingRuleLabelFingerprint(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRuleAllPorts(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRuleNetworkTier(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRuleServiceDirectoryRegistrations(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"namespace": flattenComputeForwardingRuleServiceDirectoryRegistrationsNamespace(original["namespace"], d, config),
			"service":   flattenComputeForwardingRuleServiceDirectoryRegistrationsService(original["service"], d, config),
		})
	}
	return transformed
}
func flattenComputeForwardingRuleServiceDirectoryRegistrationsNamespace(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRuleServiceDirectoryRegistrationsService(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRuleServiceLabel(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRuleServiceName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRuleSourceIpRanges(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRuleBaseForwardingRule(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRuleAllowPscGlobalAccess(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRuleIpVersion(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRuleTerraformLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}

	transformed := make(map[string]interface{})
	if l, ok := d.GetOkExists("terraform_labels"); ok {
		for k := range l.(map[string]interface{}) {
			transformed[k] = v.(map[string]interface{})[k]
		}
	}

	return transformed
}

func flattenComputeForwardingRuleEffectiveLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeForwardingRuleRegion(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.NameFromSelfLinkStateFunc(v)
}
