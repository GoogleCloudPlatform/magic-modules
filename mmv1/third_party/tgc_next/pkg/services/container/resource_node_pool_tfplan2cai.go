package container

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tfplan2cai/converters/cai"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tpgresource"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/transport"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/container/v1"
)

func ContainerNodePoolTfplan2caiConverter() cai.Tfplan2caiConverter {
	return cai.Tfplan2caiConverter{
		Convert: GetContainerNodePoolCaiObject,
	}
}

func GetContainerNodePoolCaiObject(d tpgresource.TerraformResourceData, config *transport.Config) ([]caiasset.Asset, error) {
	name, err := cai.AssetName(d, config, "//container.googleapis.com/projects/{{project}}/locations/{{location}}/clusters/{{cluster}}/nodePools/{{name}}")
	if v, ok := d.GetOk("location"); ok && tpgresource.IsZone(v.(string)) {
		name = strings.Replace(name, "/locations/", "/zones/", 1)
	}
	if err != nil {
		return []caiasset.Asset{}, err
	}
	if obj, err := GetContainerNodePoolApiObject(d, config); err == nil {
		return []caiasset.Asset{{
			Name: name,
			Type: ContainerNodePoolAssetType,
			Resource: &caiasset.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/container/v1/rest",
				DiscoveryName:        "NodePool",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []caiasset.Asset{}, err
	}
}

func GetContainerNodePoolApiObject(d tpgresource.TerraformResourceData, config *transport.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	nameProp, err := expandContainerNodePoolName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	configProp, err := expandContainerNodePoolNodeConfig(d.Get("node_config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("node_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(configProp)) && (ok || !reflect.DeepEqual(v, configProp)) {
		obj["config"] = configProp
	}
	initialNodeCountProp, err := expandContainerNodePoolInitialNodeCount(d.Get("initial_node_count"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("initial_node_count"); !tpgresource.IsEmptyValue(reflect.ValueOf(initialNodeCountProp)) && (ok || !reflect.DeepEqual(v, initialNodeCountProp)) {
		obj["initialNodeCount"] = initialNodeCountProp
	}
	versionProp, err := expandContainerNodePoolVersion(d.Get("version"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("version"); !tpgresource.IsEmptyValue(reflect.ValueOf(versionProp)) && (ok || !reflect.DeepEqual(v, versionProp)) {
		obj["version"] = versionProp
	}
	autoscalingProp, err := expandContainerNodePoolAutoscaling(d.Get("autoscaling"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("autoscaling"); !tpgresource.IsEmptyValue(reflect.ValueOf(autoscalingProp)) && (ok || !reflect.DeepEqual(v, autoscalingProp)) {
		obj["autoscaling"] = autoscalingProp
	}
	managementProp, err := expandContainerNodePoolManagement(d.Get("management"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("management"); !tpgresource.IsEmptyValue(reflect.ValueOf(managementProp)) && (ok || !reflect.DeepEqual(v, managementProp)) {
		obj["management"] = managementProp
	}
	clusterProp, err := expandContainerNodePoolCluster(d.Get("cluster"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("cluster"); !tpgresource.IsEmptyValue(reflect.ValueOf(clusterProp)) && (ok || !reflect.DeepEqual(v, clusterProp)) {
		obj["cluster"] = clusterProp
	}
	locationProp, err := expandContainerNodePoolLocation(d.Get("location"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("location"); !tpgresource.IsEmptyValue(reflect.ValueOf(locationProp)) && (ok || !reflect.DeepEqual(v, locationProp)) {
		obj["location"] = locationProp
	}

	return obj, nil
}

func expandContainerNodePoolName(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedMachineType, err := expandContainerNodePoolNodeConfigMachineType(original["machine_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMachineType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["machineType"] = transformedMachineType
	}

	transformedDiskSizeGb, err := expandContainerNodePoolNodeConfigDiskSizeGb(original["disk_size_gb"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDiskSizeGb); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["diskSizeGb"] = transformedDiskSizeGb
	}

	transformedServiceAccount, err := expandContainerNodePoolNodeConfigServiceAccount(original["service_account"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedServiceAccount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["serviceAccount"] = transformedServiceAccount
	}

	transformedMetadata, err := expandContainerNodePoolNodeConfigMetadata(original["metadata"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMetadata); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["metadata"] = transformedMetadata
	}

	transformedImageType, err := expandContainerNodePoolNodeConfigImageType(original["image_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedImageType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["imageType"] = transformedImageType
	}

	transformedLabels, err := expandContainerNodePoolNodeConfigLabels(original["labels"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedLabels); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["labels"] = transformedLabels
	}

	transformedLocalSsdCount, err := expandContainerNodePoolNodeConfigLocalSsdCount(original["local_ssd_count"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedLocalSsdCount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["localSsdCount"] = transformedLocalSsdCount
	}

	transformedTags, err := expandContainerNodePoolNodeConfigTags(original["tags"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedTags); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["tags"] = transformedTags
	}

	transformedPreemptible, err := expandContainerNodePoolNodeConfigPreemptible(original["preemptible"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPreemptible); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["preemptible"] = transformedPreemptible
	}

	transformedGuestAccelerator, err := expandContainerNodePoolNodeConfigGuestAccelerator(original["guest_accelerator"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGuestAccelerator); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["accelerators"] = transformedGuestAccelerator
	}

	transformedDiskType, err := expandContainerNodePoolNodeConfigDiskType(original["disk_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDiskType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["diskType"] = transformedDiskType
	}

	transformedMinCpuPlatform, err := expandContainerNodePoolNodeConfigMinCpuPlatform(original["min_cpu_platform"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMinCpuPlatform); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["minCpuPlatform"] = transformedMinCpuPlatform
	}

	transformedTaint, err := expandContainerNodePoolNodeConfigTaint(original["taint"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedTaint); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["taints"] = transformedTaint
	}

	return transformed, nil
}

func expandContainerNodePoolNodeConfigMachineType(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigDiskSizeGb(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigServiceAccount(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigMetadata(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandContainerNodePoolNodeConfigImageType(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandContainerNodePoolNodeConfigLocalSsdCount(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigTags(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigPreemptible(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigGuestAccelerator(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedCount, err := expandContainerNodePoolNodeConfigGuestAcceleratorCount(original["count"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedCount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["acceleratorCount"] = transformedCount
		}

		transformedType, err := expandContainerNodePoolNodeConfigGuestAcceleratorType(original["type"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["acceleratorType"] = transformedType
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandContainerNodePoolNodeConfigGuestAcceleratorCount(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	return fmt.Sprintf("%d", v.(int)), nil
}

func expandContainerNodePoolNodeConfigGuestAcceleratorType(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigDiskType(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigMinCpuPlatform(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigTaint(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedKey, err := expandContainerNodePoolNodeConfigTaintKey(original["key"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedKey); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["key"] = transformedKey
		}

		transformedValue, err := expandContainerNodePoolNodeConfigTaintValue(original["value"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedValue); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["value"] = transformedValue
		}

		transformedEffect, err := expandContainerNodePoolNodeConfigTaintEffect(original["effect"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedEffect); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["effect"] = transformedEffect
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandContainerNodePoolNodeConfigTaintKey(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigTaintValue(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolNodeConfigTaintEffect(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolInitialNodeCount(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolVersion(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolAutoscaling(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedEnabled, err := expandContainerNodePoolAutoscalingEnabled(original["enabled"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEnabled); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["enabled"] = transformedEnabled
	}

	transformedMinNodeCount, err := expandContainerNodePoolAutoscalingMinNodeCount(original["min_node_count"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMinNodeCount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["minNodeCount"] = transformedMinNodeCount
	}

	transformedMaxNodeCount, err := expandContainerNodePoolAutoscalingMaxNodeCount(original["max_node_count"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMaxNodeCount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["maxNodeCount"] = transformedMaxNodeCount
	}

	return transformed, nil
}

func expandContainerNodePoolAutoscalingEnabled(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolAutoscalingMinNodeCount(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolAutoscalingMaxNodeCount(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolManagement(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedAutoUpgrade, err := expandContainerNodePoolManagementAutoUpgrade(original["auto_upgrade"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAutoUpgrade); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["autoUpgrade"] = transformedAutoUpgrade
	}

	transformedAutoRepair, err := expandContainerNodePoolManagementAutoRepair(original["auto_repair"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAutoRepair); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["autoRepair"] = transformedAutoRepair
	}

	return transformed, nil
}

func expandContainerNodePoolManagementAutoUpgrade(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolManagementAutoRepair(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	return v, nil
}

func expandContainerNodePoolCluster(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	f, err := tpgresource.ParseGlobalFieldValue("clusters", v.(string), "project", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for cluster: %s", err)
	}
	return f.RelativeLink(), nil
}

func expandContainerNodePoolLocation(v interface{}, d tpgresource.TerraformResourceData, config *transport.Config) (interface{}, error) {
	return v, nil
}

func expandNodePool(d tpgresource.TerraformResourceData, prefix string) (*container.NodePool, error) {
	var name string
	if v, ok := d.GetOk(prefix + "name"); ok {
		if _, ok := d.GetOk(prefix + "name_prefix"); ok {
			return nil, fmt.Errorf("Cannot specify both name and name_prefix for a node_pool")
		}
		name = v.(string)
	} else if v, ok := d.GetOk(prefix + "name_prefix"); ok {
		name = id.PrefixedUniqueId(v.(string))
	} else {
		name = id.UniqueId()
	}

	nodeCount := 0
	if initialNodeCount, ok := d.GetOk(prefix + "initial_node_count"); ok {
		nodeCount = initialNodeCount.(int)
	}
	if nc, ok := d.GetOk(prefix + "node_count"); ok {
		if nodeCount != 0 {
			return nil, fmt.Errorf("Cannot set both initial_node_count and node_count on node pool %s", name)
		}
		nodeCount = nc.(int)
	}

	var locations []string
	if v, ok := d.GetOk("node_locations"); ok && v.(*schema.Set).Len() > 0 {
		locations = tpgresource.ConvertStringSet(v.(*schema.Set))
	}

	np := &container.NodePool{
		Name:             name,
		InitialNodeCount: int64(nodeCount),
		Config:           expandNodeConfig(d, prefix, d.Get(prefix+"node_config")),
		Locations:        locations,
		Version:          d.Get(prefix + "version").(string),
		NetworkConfig:    expandNodeNetworkConfig(d.Get(prefix + "network_config")),
	}

	if v, ok := d.GetOk(prefix + "autoscaling"); ok {
		if autoscaling, ok := v.([]interface{})[0].(map[string]interface{}); ok {
			np.Autoscaling = &container.NodePoolAutoscaling{
				Enabled:           true,
				MinNodeCount:      int64(autoscaling["min_node_count"].(int)),
				MaxNodeCount:      int64(autoscaling["max_node_count"].(int)),
				TotalMinNodeCount: int64(autoscaling["total_min_node_count"].(int)),
				TotalMaxNodeCount: int64(autoscaling["total_max_node_count"].(int)),
				LocationPolicy:    autoscaling["location_policy"].(string),
				ForceSendFields:   []string{"MinNodeCount", "MaxNodeCount", "TotalMinNodeCount", "TotalMaxNodeCount"},
			}
		}
	}

	if v, ok := d.GetOk(prefix + "placement_policy"); ok {
		if v.([]interface{}) != nil && v.([]interface{})[0] != nil {
			placement_policy := v.([]interface{})[0].(map[string]interface{})
			np.PlacementPolicy = &container.PlacementPolicy{
				Type:        placement_policy["type"].(string),
				PolicyName:  placement_policy["policy_name"].(string),
				TpuTopology: placement_policy["tpu_topology"].(string),
			}
		}
	}

	if v, ok := d.GetOk(prefix + "queued_provisioning"); ok {
		if v.([]interface{}) != nil && v.([]interface{})[0] != nil {
			queued_provisioning := v.([]interface{})[0].(map[string]interface{})
			np.QueuedProvisioning = &container.QueuedProvisioning{
				Enabled: queued_provisioning["enabled"].(bool),
			}
		}
	}

	if v, ok := d.GetOk(prefix + "max_pods_per_node"); ok {
		np.MaxPodsConstraint = &container.MaxPodsConstraint{
			MaxPodsPerNode: int64(v.(int)),
		}
	}

	if v, ok := d.GetOk(prefix + "management"); ok {
		managementConfig := v.([]interface{})[0].(map[string]interface{})
		np.Management = &container.NodeManagement{}

		if v, ok := managementConfig["auto_repair"]; ok {
			np.Management.AutoRepair = v.(bool)
		}

		if v, ok := managementConfig["auto_upgrade"]; ok {
			np.Management.AutoUpgrade = v.(bool)
		}
	}

	if v, ok := d.GetOk(prefix + "upgrade_settings"); ok {
		upgradeSettingsConfig := v.([]interface{})[0].(map[string]interface{})
		np.UpgradeSettings = &container.UpgradeSettings{}

		if v, ok := upgradeSettingsConfig["strategy"]; ok {
			np.UpgradeSettings.Strategy = v.(string)
		}

		if d.HasChange(prefix + "upgrade_settings.0.max_surge") {
			if np.UpgradeSettings.Strategy != "SURGE" {
				return nil, fmt.Errorf("Surge upgrade settings may not be changed when surge strategy is not enabled")
			}
			if v, ok := upgradeSettingsConfig["max_surge"]; ok {
				np.UpgradeSettings.MaxSurge = int64(v.(int))
			}
		}

		if d.HasChange(prefix + "upgrade_settings.0.max_unavailable") {
			if np.UpgradeSettings.Strategy != "SURGE" {
				return nil, fmt.Errorf("Surge upgrade settings may not be changed when surge strategy is not enabled")
			}
			if v, ok := upgradeSettingsConfig["max_unavailable"]; ok {
				np.UpgradeSettings.MaxUnavailable = int64(v.(int))
			}
		}

		if v, ok := upgradeSettingsConfig["blue_green_settings"]; ok && len(v.([]interface{})) > 0 {
			blueGreenSettingsConfig := v.([]interface{})[0].(map[string]interface{})
			np.UpgradeSettings.BlueGreenSettings = &container.BlueGreenSettings{}

			if np.UpgradeSettings.Strategy != "BLUE_GREEN" {
				return nil, fmt.Errorf("Blue-green upgrade settings may not be changed when blue-green strategy is not enabled")
			}

			if v, ok := blueGreenSettingsConfig["node_pool_soak_duration"]; ok {
				np.UpgradeSettings.BlueGreenSettings.NodePoolSoakDuration = v.(string)
			}

			if v, ok := blueGreenSettingsConfig["standard_rollout_policy"]; ok && len(v.([]interface{})) > 0 {
				standardRolloutPolicyConfig := v.([]interface{})[0].(map[string]interface{})
				standardRolloutPolicy := &container.StandardRolloutPolicy{}

				if v, ok := standardRolloutPolicyConfig["batch_soak_duration"]; ok {
					standardRolloutPolicy.BatchSoakDuration = v.(string)
				}
				if v, ok := standardRolloutPolicyConfig["batch_node_count"]; ok {
					standardRolloutPolicy.BatchNodeCount = int64(v.(int))
				}
				if v, ok := standardRolloutPolicyConfig["batch_percentage"]; ok {
					standardRolloutPolicy.BatchPercentage = v.(float64)
				}

				np.UpgradeSettings.BlueGreenSettings.StandardRolloutPolicy = standardRolloutPolicy
			}

		}
	}

	if v, ok := d.GetOk(prefix + "node_drain_config"); ok {
		nodeDrainConfig := v.([]interface{})[0].(map[string]interface{})
		np.NodeDrainConfig = &container.NodeDrainConfig{}

		if v, ok := nodeDrainConfig["respect_pdb_during_node_pool_deletion"]; ok {
			np.NodeDrainConfig.RespectPdbDuringNodePoolDeletion = v.(bool)
		}
	}

	return np, nil
}

func expandNodeNetworkConfig(v interface{}) *container.NodeNetworkConfig {
	networkNodeConfigs := v.([]interface{})

	nnc := &container.NodeNetworkConfig{}

	if len(networkNodeConfigs) == 0 {
		return nnc
	}

	networkNodeConfig := networkNodeConfigs[0].(map[string]interface{})

	if v, ok := networkNodeConfig["create_pod_range"]; ok {
		nnc.CreatePodRange = v.(bool)
	}

	if v, ok := networkNodeConfig["pod_range"]; ok {
		nnc.PodRange = v.(string)
	}

	if v, ok := networkNodeConfig["pod_ipv4_cidr_block"]; ok {
		nnc.PodIpv4CidrBlock = v.(string)
	}

	if v, ok := networkNodeConfig["enable_private_nodes"]; ok {
		nnc.EnablePrivateNodes = v.(bool)
		nnc.ForceSendFields = []string{"EnablePrivateNodes"}
	}

	if v, ok := networkNodeConfig["additional_node_network_configs"]; ok && len(v.([]interface{})) > 0 {
		node_network_configs := v.([]interface{})
		nodeNetworkConfigs := make([]*container.AdditionalNodeNetworkConfig, 0, len(node_network_configs))
		for _, raw := range node_network_configs {
			data := raw.(map[string]interface{})
			networkConfig := &container.AdditionalNodeNetworkConfig{
				Network:    data["network"].(string),
				Subnetwork: data["subnetwork"].(string),
			}
			nodeNetworkConfigs = append(nodeNetworkConfigs, networkConfig)
		}
		nnc.AdditionalNodeNetworkConfigs = nodeNetworkConfigs
	}

	if v, ok := networkNodeConfig["additional_pod_network_configs"]; ok && len(v.([]interface{})) > 0 {
		pod_network_configs := v.([]interface{})
		podNetworkConfigs := make([]*container.AdditionalPodNetworkConfig, 0, len(pod_network_configs))
		for _, raw := range pod_network_configs {
			data := raw.(map[string]interface{})
			podnetworkConfig := &container.AdditionalPodNetworkConfig{
				Subnetwork:        data["subnetwork"].(string),
				SecondaryPodRange: data["secondary_pod_range"].(string),
				MaxPodsPerNode: &container.MaxPodsConstraint{
					MaxPodsPerNode: int64(data["max_pods_per_node"].(int)),
				},
			}
			podNetworkConfigs = append(podNetworkConfigs, podnetworkConfig)
		}
		nnc.AdditionalPodNetworkConfigs = podNetworkConfigs
	}

	nnc.PodCidrOverprovisionConfig = expandPodCidrOverprovisionConfig(networkNodeConfig["pod_cidr_overprovision_config"])

	if v, ok := networkNodeConfig["network_performance_config"]; ok && len(v.([]interface{})) > 0 {
		nnc.NetworkPerformanceConfig = &container.NetworkPerformanceConfig{}
		network_performance_config := v.([]interface{})[0].(map[string]interface{})
		if total_egress_bandwidth_tier, ok := network_performance_config["total_egress_bandwidth_tier"]; ok {
			nnc.NetworkPerformanceConfig.TotalEgressBandwidthTier = total_egress_bandwidth_tier.(string)
		}
	}

	return nnc
}
