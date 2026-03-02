package container

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/cai2hcl/converters/utils"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/cai2hcl/models"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/transport"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ContainerNodePoolCai2hclConverter for container node pool resource.
type ContainerNodePoolCai2hclConverter struct {
	name   string
	schema map[string]*schema.Schema
}

// NewContainerNodePoolCai2hclConverter returns an HCL converter for container node pool.
func NewContainerNodePoolCai2hclConverter(provider *schema.Provider) models.Cai2hclConverter {
	schema := provider.ResourcesMap["google_container_node_pool"].Schema

	return &ContainerNodePoolCai2hclConverter{
		name:   "google_container_node_pool",
		schema: schema,
	}
}

// Convert converts asset resource data.
func (c *ContainerNodePoolCai2hclConverter) Convert(asset caiasset.Asset) ([]*models.TerraformResourceBlock, error) {
	if asset.Resource == nil || asset.Resource.Data == nil {
		return nil, fmt.Errorf("asset resource data is nil")
	}

	block, err := c.convertResourceData(asset)
	if err != nil {
		return nil, err
	}
	return []*models.TerraformResourceBlock{block}, nil
}

func (c *ContainerNodePoolCai2hclConverter) convertResourceData(asset caiasset.Asset) (*models.TerraformResourceBlock, error) {
	if asset.Resource == nil || asset.Resource.Data == nil {
		return nil, fmt.Errorf("asset resource data is nil")
	}

	config := transport.NewConfig()

	// This is a fake resource used to get fake d
	// d.Get will return empty map, instead of nil
	fakeResource := &schema.Resource{
		Schema: c.schema,
	}
	d := fakeResource.TestResourceData()

	hclData := make(map[string]interface{})

	outputFields := map[string]struct{}{}
	an := strings.Replace(asset.Name, "/zones/", "/locations/", 1)
	utils.ParseUrlParamValuesFromAssetName(an, "//container.googleapis.com/projects/{{project}}/locations/{{location}}/clusters/{{cluster}}/nodePools/{{name}}", outputFields, hclData)

	npMap, err := flattenNodePool(d, config, asset.Resource.Data, "")
	if err != nil {
		return nil, err
	}

	for k, v := range npMap {
		hclData[k] = v
	}

	ctyVal, err := utils.MapToCtyValWithSchema(hclData, c.schema)
	if err != nil {
		return nil, err
	}
	return &models.TerraformResourceBlock{
		Labels: []string{c.name, asset.Resource.Data["name"].(string)},
		Value:  ctyVal,
	}, nil
}

func flattenNodePoolStandardRolloutPolicy(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	rp, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"batch_node_count":    rp["batchNodeCount"],
		"batch_percentage":    rp["batchPercentage"],
		"batch_soak_duration": rp["batchSoakDuration"],
	}
	return []map[string]interface{}{transformed}
}

func flattenNodePoolBlueGreenSettings(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	bg, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{
		"node_pool_soak_duration": bg["nodePoolSoakDuration"],
		"standard_rollout_policy": flattenNodePoolStandardRolloutPolicy(bg["standardRolloutPolicy"]),
	}
	return []map[string]interface{}{transformed}
}

func flattenNodePoolUpgradeSettings(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	us, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	upgradeSettings := make(map[string]interface{})

	upgradeSettings["blue_green_settings"] = flattenNodePoolBlueGreenSettings(us["blueGreenSettings"])
	upgradeSettings["max_surge"] = us["maxSurge"]
	upgradeSettings["max_unavailable"] = us["maxUnavailable"]

	// "SHORT_LIVED" strategy is not supported by the Terraform provider yet.
	if strategy, ok := us["strategy"].(string); ok && strategy != "SHORT_LIVED" {
		upgradeSettings["strategy"] = strategy
	}

	return []map[string]interface{}{upgradeSettings}
}

func flattenNodePoolNodeDrainConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	ndc, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	nodeDrainConfig := make(map[string]interface{})
	nodeDrainConfig["respect_pdb_during_node_pool_deletion"] = ndc["respectPdbDuringNodePoolDeletion"]

	return []map[string]interface{}{nodeDrainConfig}
}

func flattenNodePool(d *schema.ResourceData, config *transport.Config, np map[string]interface{}, prefix string) (map[string]interface{}, error) {
	// Node pools don't expose the current node count in their API, so read the
	// instance groups instead. They should all have the same size, but in case a resize
	// failed or something else strange happened, we'll just use the average size.
	size := 0
	igmUrls := []string{}

	if v, ok := np["instanceGroupUrls"].([]interface{}); ok {
		for _, url := range v {
			if urlStr, ok := url.(string); ok {
				// retrieve instance group manager (InstanceGroupUrls are actually URLs for InstanceGroupManagers)
				matches := instanceGroupManagerURL.FindStringSubmatch(urlStr)
				if len(matches) < 4 {
					return nil, fmt.Errorf("Error reading instance group manage URL '%q'", urlStr)
				}
				if strings.HasPrefix("gk3", matches[3]) {
					// IGM is autopilot so we know it will not be found, skip it
					continue
				}
				igmUrls = append(igmUrls, urlStr)
			}
		}
	}

	nodeCount := 0
	if len(igmUrls) > 0 {
		nodeCount = size / len(igmUrls)
	}
	nodePool := map[string]interface{}{
		"name":               np["name"],
		"initial_node_count": np["initialNodeCount"],
		"node_locations":     np["locations"],
		"node_config":        flattenNodeConfig(np["config"], d.Get(prefix+"node_config")),
		"version":            np["version"],
		"network_config":     flattenNodeNetworkConfig(np["networkConfig"], d, prefix),
	}

	if nodeCount > 0 {
		nodePool["node_count"] = nodeCount
	}

	if v, ok := np["autoscaling"].(map[string]interface{}); ok {
		if enabled, ok := v["enabled"].(bool); ok && enabled {
			nodePool["autoscaling"] = []map[string]interface{}{
				{
					"min_node_count":       v["minNodeCount"],
					"max_node_count":       v["maxNodeCount"],
					"total_min_node_count": v["totalMinNodeCount"],
					"total_max_node_count": v["totalMaxNodeCount"],
					"location_policy":      v["locationPolicy"],
				},
			}
		} else {
			nodePool["autoscaling"] = []map[string]interface{}{}
		}
	}

	if v, ok := np["placementPolicy"].(map[string]interface{}); ok {
		nodePool["placement_policy"] = []map[string]interface{}{
			{
				"type":         v["type"],
				"policy_name":  v["policyName"],
				"tpu_topology": v["tpuTopology"],
			},
		}
	}

	if v, ok := np["queuedProvisioning"].(map[string]interface{}); ok {
		nodePool["queued_provisioning"] = []map[string]interface{}{
			{
				"enabled": v["enabled"],
			},
		}
	}

	if v, ok := np["maxPodsConstraint"].(map[string]interface{}); ok {
		nodePool["max_pods_per_node"] = v["maxPodsPerNode"]
	}

	if v, ok := np["management"].(map[string]interface{}); ok {
		nodePool["management"] = []map[string]interface{}{
			{
				"auto_repair":  v["autoRepair"],
				"auto_upgrade": v["autoUpgrade"],
			},
		}
	}

	nodePool["upgrade_settings"] = flattenNodePoolUpgradeSettings(np["upgradeSettings"])

	nodePool["node_drain_config"] = flattenNodePoolNodeDrainConfig(np["nodeDrainConfig"])

	return nodePool, nil
}

func flattenNodeNetworkConfig(c interface{}, d *schema.ResourceData, prefix string) []map[string]interface{} {
	if c == nil {
		return nil
	}

	if config, ok := c.(map[string]interface{}); ok {
		transformed := map[string]interface{}{
			// TODO: investigate why create_pod_range is not returned by the API
			// "create_pod_range": d.Get(prefix + "network_config.0.create_pod_range"), // API doesn't return this value so we set the old one. Field is ForceNew + Required
			"pod_ipv4_cidr_block":             config["podIpv4CidrBlock"],
			"pod_range":                       config["podRange"],
			"enable_private_nodes":            config["enablePrivateNodes"],
			"pod_cidr_overprovision_config":   flattenPodCidrOverprovisionConfig(config["podCidrOverprovisionConfig"]),
			"network_performance_config":      flattenNodeNetworkPerformanceConfig(config["networkPerformanceConfig"]),
			"additional_node_network_configs": flattenAdditionalNodeNetworkConfig(config["additionalNodeNetworkConfigs"]),
			"additional_pod_network_configs":  flattenAdditionalPodNetworkConfig(config["additionalPodNetworkConfigs"]),
		}
		return []map[string]interface{}{transformed}
	}
	return nil
}

func flattenNodeNetworkPerformanceConfig(c interface{}) []map[string]interface{} {
	if c == nil {
		return nil
	}
	if config, ok := c.(map[string]interface{}); ok {
		transformed := map[string]interface{}{
			"total_egress_bandwidth_tier": config["totalEgressBandwidthTier"],
		}
		return []map[string]interface{}{transformed}
	}
	return nil
}

func flattenAdditionalNodeNetworkConfig(c interface{}) []map[string]interface{} {
	if c == nil {
		return nil
	}
	result := []map[string]interface{}{}
	if configs, ok := c.([]interface{}); ok {
		for _, v := range configs {
			if config, ok := v.(map[string]interface{}); ok {
				transformed := map[string]interface{}{
					"network":    config["network"],
					"subnetwork": config["subnetwork"],
				}
				result = append(result, transformed)
			}
		}
	}
	return result
}

func flattenAdditionalPodNetworkConfig(c interface{}) []map[string]interface{} {
	if c == nil {
		return nil
	}
	result := []map[string]interface{}{}
	if configs, ok := c.([]interface{}); ok {
		for _, v := range configs {
			if config, ok := v.(map[string]interface{}); ok {
				var maxPodsPerNode interface{}
				if mp, ok := config["maxPodsPerNode"].(map[string]interface{}); ok {
					maxPodsPerNode = mp["maxPodsPerNode"]
				}
				transformed := map[string]interface{}{
					"subnetwork":          config["subnetwork"],
					"secondary_pod_range": config["secondaryPodRange"],
					"max_pods_per_node":   maxPodsPerNode,
				}
				result = append(result, transformed)
			}
		}
	}
	return result
}
