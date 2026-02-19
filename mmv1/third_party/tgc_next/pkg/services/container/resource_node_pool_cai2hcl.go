package container

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/cai2hcl/converters/utils"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/cai2hcl/models"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tpgresource"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/transport"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/container/v1"
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

	var nodePool *container.NodePool
	if err := utils.DecodeJSON(asset.Resource.Data, &nodePool); err != nil {
		return nil, err
	}

	hclData := make(map[string]interface{})

	npMap, err := flattenNodePool(d, config, nodePool, "")
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
		Labels: []string{c.name, nodePool.Name},
		Value:  ctyVal,
	}, nil
}

func flattenNodePoolStandardRolloutPolicy(rp *container.StandardRolloutPolicy) []map[string]interface{} {
	if rp == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"batch_node_count":    rp.BatchNodeCount,
			"batch_percentage":    rp.BatchPercentage,
			"batch_soak_duration": rp.BatchSoakDuration,
		},
	}
}

func flattenNodePoolBlueGreenSettings(bg *container.BlueGreenSettings) []map[string]interface{} {
	if bg == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"node_pool_soak_duration": bg.NodePoolSoakDuration,
			"standard_rollout_policy": flattenNodePoolStandardRolloutPolicy(bg.StandardRolloutPolicy),
		},
	}
}

func flattenNodePoolUpgradeSettings(us *container.UpgradeSettings) []map[string]interface{} {
	if us == nil {
		return nil
	}

	upgradeSettings := make(map[string]interface{})

	upgradeSettings["blue_green_settings"] = flattenNodePoolBlueGreenSettings(us.BlueGreenSettings)
	upgradeSettings["max_surge"] = us.MaxSurge
	upgradeSettings["max_unavailable"] = us.MaxUnavailable

	upgradeSettings["strategy"] = us.Strategy
	return []map[string]interface{}{upgradeSettings}
}

func flattenNodePoolNodeDrainConfig(ndc *container.NodeDrainConfig) []map[string]interface{} {
	if ndc == nil {
		return nil
	}

	nodeDrainConfig := make(map[string]interface{})

	nodeDrainConfig["respect_pdb_during_node_pool_deletion"] = ndc.RespectPdbDuringNodePoolDeletion
	return []map[string]interface{}{nodeDrainConfig}
}

func flattenNodePool(d *schema.ResourceData, config *transport.Config, np *container.NodePool, prefix string) (map[string]interface{}, error) {
	// Node pools don't expose the current node count in their API, so read the
	// instance groups instead. They should all have the same size, but in case a resize
	// failed or something else strange happened, we'll just use the average size.
	size := 0
	igmUrls := []string{}
	managedIgmUrls := []string{}
	for _, url := range np.InstanceGroupUrls {
		// retrieve instance group manager (InstanceGroupUrls are actually URLs for InstanceGroupManagers)
		matches := instanceGroupManagerURL.FindStringSubmatch(url)
		if len(matches) < 4 {
			return nil, fmt.Errorf("Error reading instance group manage URL '%q'", url)
		}
		if strings.HasPrefix("gk3", matches[3]) {
			// IGM is autopilot so we know it will not be found, skip it
			continue
		}

		// TODO: get igms
	}
	nodeCount := 0
	if len(igmUrls) > 0 {
		nodeCount = size / len(igmUrls)
	}
	nodePool := map[string]interface{}{
		"name":                        np.Name,
		"name_prefix":                 d.Get(prefix + "name_prefix"),
		"initial_node_count":          np.InitialNodeCount,
		"node_locations":              schema.NewSet(schema.HashString, tpgresource.ConvertStringArrToInterface(np.Locations)),
		"node_count":                  nodeCount,
		"node_config":                 flattenNodeConfig(np.Config, d.Get(prefix+"node_config")),
		"instance_group_urls":         igmUrls,
		"managed_instance_group_urls": managedIgmUrls,
		"version":                     np.Version,
		"network_config":              flattenNodeNetworkConfig(np.NetworkConfig, d, prefix),
	}

	if np.Autoscaling != nil {
		if np.Autoscaling.Enabled {
			nodePool["autoscaling"] = []map[string]interface{}{
				{
					"min_node_count":       np.Autoscaling.MinNodeCount,
					"max_node_count":       np.Autoscaling.MaxNodeCount,
					"total_min_node_count": np.Autoscaling.TotalMinNodeCount,
					"total_max_node_count": np.Autoscaling.TotalMaxNodeCount,
					"location_policy":      np.Autoscaling.LocationPolicy,
				},
			}
		} else {
			nodePool["autoscaling"] = []map[string]interface{}{}
		}
	}

	if np.PlacementPolicy != nil {
		nodePool["placement_policy"] = []map[string]interface{}{
			{
				"type":         np.PlacementPolicy.Type,
				"policy_name":  np.PlacementPolicy.PolicyName,
				"tpu_topology": np.PlacementPolicy.TpuTopology,
			},
		}
	}

	if np.QueuedProvisioning != nil {
		nodePool["queued_provisioning"] = []map[string]interface{}{
			{
				"enabled": np.QueuedProvisioning.Enabled,
			},
		}
	}

	if np.MaxPodsConstraint != nil {
		nodePool["max_pods_per_node"] = np.MaxPodsConstraint.MaxPodsPerNode
	}

	if np.Management != nil {
		nodePool["management"] = []map[string]interface{}{
			{
				"auto_repair":  np.Management.AutoRepair,
				"auto_upgrade": np.Management.AutoUpgrade,
			},
		}
	}

	if np.UpgradeSettings != nil {
		nodePool["upgrade_settings"] = flattenNodePoolUpgradeSettings(np.UpgradeSettings)
	} else {
		delete(nodePool, "upgrade_settings")
	}

	if np.NodeDrainConfig != nil {
		nodePool["node_drain_config"] = flattenNodePoolNodeDrainConfig(np.NodeDrainConfig)
	}

	return nodePool, nil
}

func flattenNodeNetworkConfig(c *container.NodeNetworkConfig, d *schema.ResourceData, prefix string) []map[string]interface{} {
	result := []map[string]interface{}{}
	if c != nil {
		result = append(result, map[string]interface{}{
			// TODO: investigate why create_pod_range is not returned by the API
			// "create_pod_range": d.Get(prefix + "network_config.0.create_pod_range"), // API doesn't return this value so we set the old one. Field is ForceNew + Required
			"pod_ipv4_cidr_block":             c.PodIpv4CidrBlock,
			"pod_range":                       c.PodRange,
			"enable_private_nodes":            c.EnablePrivateNodes,
			"pod_cidr_overprovision_config":   flattenPodCidrOverprovisionConfig(c.PodCidrOverprovisionConfig),
			"network_performance_config":      flattenNodeNetworkPerformanceConfig(c.NetworkPerformanceConfig),
			"additional_node_network_configs": flattenAdditionalNodeNetworkConfig(c.AdditionalNodeNetworkConfigs),
			"additional_pod_network_configs":  flattenAdditionalPodNetworkConfig(c.AdditionalPodNetworkConfigs),
			"subnetwork":                      c.Subnetwork,
		})
	}
	return result
}

func flattenNodeNetworkPerformanceConfig(c *container.NetworkPerformanceConfig) []map[string]interface{} {
	result := []map[string]interface{}{}
	if c != nil {
		result = append(result, map[string]interface{}{
			"total_egress_bandwidth_tier": c.TotalEgressBandwidthTier,
		})
	}
	return result
}

func flattenAdditionalNodeNetworkConfig(c []*container.AdditionalNodeNetworkConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}

	result := []map[string]interface{}{}
	for _, nodeNetworkConfig := range c {
		result = append(result, map[string]interface{}{
			"network":    nodeNetworkConfig.Network,
			"subnetwork": nodeNetworkConfig.Subnetwork,
		})
	}
	return result
}

func flattenAdditionalPodNetworkConfig(c []*container.AdditionalPodNetworkConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}

	result := []map[string]interface{}{}
	for _, podNetworkConfig := range c {
		result = append(result, map[string]interface{}{
			"subnetwork":          podNetworkConfig.Subnetwork,
			"secondary_pod_range": podNetworkConfig.SecondaryPodRange,
			"max_pods_per_node":   podNetworkConfig.MaxPodsPerNode.MaxPodsPerNode,
		})
	}
	return result
}
