package appengine

import (
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const AppEngineVersionAssetType string = "appengine.googleapis.com/Version"

func ResourceAppEngineStandardAppVersion() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: AppEngineVersionAssetType,
		Convert:   GetAppEngineVersionCaiObject,
	}
}

func GetAppEngineVersionCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//appengine.googleapis.com/apps/{{project}}/services/default/versions/v1")
	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetAppEngineVersionApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: AppEngineVersionAssetType,
			Resource: &cai.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/appengine/v1/rest",
				DiscoveryName:        "Version",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}

func GetAppEngineVersionApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	nameProp, err := expandAppEngineVersionName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}

	idProp, err := expandAppEngineVersionId(d.Get("version_id"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("version_id"); !tpgresource.IsEmptyValue(reflect.ValueOf(idProp)) && (ok || !reflect.DeepEqual(v, idProp)) {
		obj["id"] = idProp
	}

	automaticScalingProp, err := expandAppEngineVersionAutomaticScaling(d.Get("automatic_scaling"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("automatic_scaling"); !tpgresource.IsEmptyValue(reflect.ValueOf(automaticScalingProp)) && (ok || !reflect.DeepEqual(v, automaticScalingProp)) {
		obj["automaticScaling"] = automaticScalingProp
	}

	basicScalingProp, err := expandAppEngineVersionBasicScaling(d.Get("basic_scaling"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("basic_scaling"); !tpgresource.IsEmptyValue(reflect.ValueOf(basicScalingProp)) && (ok || !reflect.DeepEqual(v, basicScalingProp)) {
		obj["basicScaling"] = basicScalingProp
	}

	manualScalingProp, err := expandAppEngineVersionManualScaling(d.Get("manual_scaling"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("manual_scaling"); !tpgresource.IsEmptyValue(reflect.ValueOf(manualScalingProp)) && (ok || !reflect.DeepEqual(v, manualScalingProp)) {
		obj["manualScaling"] = manualScalingProp
	}

	instanceClassProp, err := expandAppEngineVersionInstanceClass(d.Get("instance_class"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("instance_class"); !tpgresource.IsEmptyValue(reflect.ValueOf(instanceClassProp)) && (ok || !reflect.DeepEqual(v, instanceClassProp)) {
		obj["instanceClass"] = instanceClassProp
	}

	zonesProp, err := expandAppEngineVersionZones(d.Get("zones"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("zones"); !tpgresource.IsEmptyValue(reflect.ValueOf(zonesProp)) && (ok || !reflect.DeepEqual(v, zonesProp)) {
		obj["zones"] = zonesProp
	}

	runtimeProp, err := expandAppEngineVersionRuntime(d.Get("runtime"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("runtime"); !tpgresource.IsEmptyValue(reflect.ValueOf(runtimeProp)) && (ok || !reflect.DeepEqual(v, runtimeProp)) {
		obj["runtime"] = runtimeProp
	}

	threadsafeProp, err := expandAppEngineVersionThreadsafe(d.Get("threadsafe"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("threadsafe"); !tpgresource.IsEmptyValue(reflect.ValueOf(threadsafeProp)) && (ok || !reflect.DeepEqual(v, threadsafeProp)) {
		obj["threadsafe"] = threadsafeProp
	}

	vmProp, err := expandAppEngineVersionVm(d.Get("vm"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("vm"); !tpgresource.IsEmptyValue(reflect.ValueOf(vmProp)) && (ok || !reflect.DeepEqual(v, vmProp)) {
		obj["vm"] = vmProp
	}

	appEngineApisProp, err := expandAppEngineVersionAppEngineApis(d.Get("app_engine_apis"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("app_engine_apis"); !tpgresource.IsEmptyValue(reflect.ValueOf(appEngineApisProp)) && (ok || !reflect.DeepEqual(v, appEngineApisProp)) {
		obj["appEngineApis"] = appEngineApisProp
	}

	envProp, err := expandAppEngineVersionEnv(d.Get("env"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("env"); !tpgresource.IsEmptyValue(reflect.ValueOf(envProp)) && (ok || !reflect.DeepEqual(v, envProp)) {
		obj["env"] = envProp
	}

	createTimeProp, err := expandAppEngineVersionCreateTime(d.Get("create"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("create"); !tpgresource.IsEmptyValue(reflect.ValueOf(createTimeProp)) && (ok || !reflect.DeepEqual(v, createTimeProp)) {
		obj["createTime"] = createTimeProp
	}

	runtimeApiVersionProp, err := expandAppEngineVersionRuntimeApiVersion(d.Get("runtime_api_version"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("runtime_api_version"); !tpgresource.IsEmptyValue(reflect.ValueOf(runtimeApiVersionProp)) && (ok || !reflect.DeepEqual(v, runtimeApiVersionProp)) {
		obj["runtimeApiVersion"] = runtimeApiVersionProp
	}

	serviceAccountProp, err := expandAppEngineVersionServiceAccount(d.Get("service_account"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("service_account"); !tpgresource.IsEmptyValue(reflect.ValueOf(serviceAccountProp)) && (ok || !reflect.DeepEqual(v, serviceAccountProp)) {
		obj["serviceAccount"] = serviceAccountProp
	}

	handlersProp, err := expandAppEngineVersionHandlers(d.Get("handlers"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("handlers"); !tpgresource.IsEmptyValue(reflect.ValueOf(handlersProp)) && (ok || !reflect.DeepEqual(v, handlersProp)) {
		obj["handlers"] = handlersProp
	}

	librariesProp, err := expandAppEngineVersionLibraries(d.Get("libraries"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("libraries"); !tpgresource.IsEmptyValue(reflect.ValueOf(librariesProp)) && (ok || !reflect.DeepEqual(v, librariesProp)) {
		obj["libraries"] = librariesProp
	}

	deploymentProp, err := expandAppEngineVersionDeployment(d.Get("deployment"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("deployment"); !tpgresource.IsEmptyValue(reflect.ValueOf(deploymentProp)) && (ok || !reflect.DeepEqual(v, deploymentProp)) {
		obj["deployment"] = deploymentProp
	}

	entrypointProp, err := expandAppEngineVersionEntrypoint(d.Get("entrypoint"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("entrypoint"); !tpgresource.IsEmptyValue(reflect.ValueOf(entrypointProp)) && (ok || !reflect.DeepEqual(v, entrypointProp)) {
		obj["entrypoint"] = entrypointProp
	}

	envVariablesProp, err := expandAppEngineVersionEnvVariables(d.Get("env_variables"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("envVariables"); !tpgresource.IsEmptyValue(reflect.ValueOf(envVariablesProp)) && (ok || !reflect.DeepEqual(v, envVariablesProp)) {
		obj["envVariables"] = envVariablesProp
	}

	vpcAccessConnectorProp, err := expandAppEngineVersionVpcAccessConnector(d.Get("vpc_access_connector"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("vpc_access_connector"); !tpgresource.IsEmptyValue(reflect.ValueOf(vpcAccessConnectorProp)) && (ok || !reflect.DeepEqual(v, vpcAccessConnectorProp)) {
		obj["vpcAccessConnector"] = vpcAccessConnectorProp
	}

	inboundServicesProp, err := expandAppEngineVersionInboundServices(d.Get("inbound_services"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("inbound_services"); !tpgresource.IsEmptyValue(reflect.ValueOf(inboundServicesProp)) && (ok || !reflect.DeepEqual(v, inboundServicesProp)) {
		obj["inboundServices"] = inboundServicesProp
	}

	projectProp, err := expandAppEngineVersionProject(d.Get("project"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("project"); !tpgresource.IsEmptyValue(reflect.ValueOf(projectProp)) && (ok || !reflect.DeepEqual(v, projectProp)) {
		obj["project"] = projectProp
	}

	noopOnDestroyProp, err := expandAppEngineVersionNoopOnDestroy(d.Get("noop_on_destroy"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("noop_on_destroy"); !tpgresource.IsEmptyValue(reflect.ValueOf(noopOnDestroyProp)) && (ok || !reflect.DeepEqual(v, noopOnDestroyProp)) {
		obj["noopOnDestroy"] = noopOnDestroyProp
	}

	deleteServiceOnDestroyProp, err := expandAppEngineVersionDeleteServiceOnDestroy(d.Get("delete_service_on_destroy"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("delete_service_on_destroy"); !tpgresource.IsEmptyValue(reflect.ValueOf(deleteServiceOnDestroyProp)) && (ok || !reflect.DeepEqual(v, deleteServiceOnDestroyProp)) {
		obj["deleteServiceOnDestroy"] = deleteServiceOnDestroyProp
	}

	serviceProp, err := expandAppEngineVersionService(d.Get("service"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("service"); !tpgresource.IsEmptyValue(reflect.ValueOf(serviceProp)) && (ok || !reflect.DeepEqual(v, serviceProp)) {
		obj["service"] = serviceProp
	}

	return obj, nil
}

func expandAppEngineVersionService(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionDeleteServiceOnDestroy(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionNoopOnDestroy(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionProject(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionInboundServices(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	return v, nil
}

func expandAppEngineVersionVpcAccessConnector(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedName, err := expandAppEngineVersionName(original["name"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["name"] = transformedName
	}

	transformedEgressSetting, err := expandAppEngineVersionEgressSetting(original["egress_setting"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEgressSetting); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["egressSetting"] = transformedEgressSetting
	}

	return transformed, nil
}

func expandAppEngineVersionEgressSetting(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionEnvVariables(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionAutomaticScaling(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedMaxConcurrentRequests, err := expandAppEngineVersionMaxConcurrentRequests(original["max_concurrent_requests"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMaxConcurrentRequests); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["maxConcurrentRequests"] = transformedMaxConcurrentRequests
	}

	transformedMaxIdleInstances, err := expandAppEngineVersionMaxIdleInstances(original["max_idle_instances"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMaxIdleInstances); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["maxIdleInstances"] = transformedMaxIdleInstances
	}

	transformedMaxPendingLatency, err := expandAppEngineVersionMaxPendingLatency(original["max_pending_latency"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMaxPendingLatency); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["maxPendingLatency"] = transformedMaxPendingLatency
	}

	transformedMinIdleInstances, err := expandAppEngineVersionMinIdleInstances(original["min_idle_instances"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMinIdleInstances); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["minIdleInstances"] = transformedMinIdleInstances
	}

	transformedMinPendingLatency, err := expandAppEngineVersionMinPendingLatency(original["min_pending_latency"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMinPendingLatency); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["minPendingLatency"] = transformedMinPendingLatency
	}

	transformedStandardSchedulerSettings, err := expandAppEngineVersionStandardSchedulerSettings(original["standard_scheduler_settings"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedStandardSchedulerSettings); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["standardSchedulerSettings"] = transformedStandardSchedulerSettings
	}

	return transformed, nil
}

func expandAppEngineVersionStandardSchedulerSettings(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedTargetCpuUtilization, err := expandAppEngineVersionTargetCpuUtilization(original["target_cpu_utilization"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedTargetCpuUtilization); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["targetCpuUtilization"] = transformedTargetCpuUtilization
	}

	transformedTargetThroughputUtilization, err := expandAppEngineVersionTargetThroughputUtilization(original["target_throughput_utilization"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedTargetThroughputUtilization); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["targetThroughputUtilization"] = transformedTargetThroughputUtilization
	}

	transformedMinInstances, err := expandAppEngineVersionMinInstances(original["min_instances"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMinInstances); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["minInstances"] = transformedMinInstances
	}

	transformedMaxInstances, err := expandAppEngineVersionMaxInstances(original["max_instances"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMaxInstances); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["maxInstances"] = transformedMaxInstances
	}

	return transformed, nil
}

func expandAppEngineVersionForwardedPorts(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionInstanceTag(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionSubnetworkName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionSessionAffinity(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionCoolDownPeriod(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionAggregationWindowLength(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionTargetUtilization(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionMaxConcurrentRequests(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionMaxIdleInstances(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionMaxTotalInstances(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionMaxPendingLatency(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionMinIdleInstances(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionMinTotalInstances(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionMinPendingLatency(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionTargetRequestCountPerSecond(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionTargetConcurrentRequests(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionTargetWriteOpsPerSecond(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionTargetReadBytesPerSecond(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionTargetWriteBytesPerSecond(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionTargetReadOpsPerSecond(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionTargetSentBytesPerSecond(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionTargetSentPacketsPerSecond(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionTargetReceivedBytesPerSecond(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionTargetReceivedPacketsPerSecond(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionTargetCpuUtilization(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionTargetThroughputUtilization(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionMinInstances(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionMaxInstances(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionBasicScaling(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedIdleTimeout, err := expandAppEngineVersionIdleTimeout(original["idle_timeout"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedIdleTimeout); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["idleTimeout"] = transformedIdleTimeout
	}

	transformedMaxInstances, err := expandAppEngineVersionMaxInstances(original["max_instances"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMaxInstances); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["maxInstances"] = transformedMaxInstances
	}

	return transformed, nil
}

func expandAppEngineVersionIdleTimeout(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionManualScaling(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedInstances, err := expandAppEngineVersionInstances(original["instances"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedInstances); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["instances"] = transformedInstances
	}

	return transformed, nil
}

func expandAppEngineVersionInstances(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionInstanceClass(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionZones(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionCpu(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionDiskGb(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionMemoryGb(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionKmsKeyReference(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionVolumeType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionSizeGb(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionRuntime(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionRuntimeChannel(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionThreadsafe(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionVm(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionOperatingSystem(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionrRuntimeVersion(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionAppEngineApis(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionEnv(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionCreateTime(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionDiskUsageBytes(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionRuntimeApiVersion(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionRuntimeMainExecutablePath(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionServiceAccount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionCreatedBy(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionUrlRegex(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionSecurityLevel(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionLogin(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionAuthFailAction(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionRedirectHttpResponseCode(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionStaticFilesHandler(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedPath, err := expandAppEngineVersionPath(original["path"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPath); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["path"] = transformedPath
	}

	transformedUploadPathRegex, err := expandAppEngineVersionUploadPathRegex(original["upload_path_regex"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedUploadPathRegex); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["uploadPathRegex"] = transformedUploadPathRegex
	}

	transformedHttpHeaders, err := expandAppEngineVersionHttpHeaders(original["http_headers"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedHttpHeaders); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["httpHeaders"] = transformedHttpHeaders
	}

	transformedMimeType, err := expandAppEngineVersionMimeType(original["mime_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMimeType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["mimeType"] = transformedMimeType
	}

	transformedExpiration, err := expandAppEngineVersionExpiration(original["expirtation"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedExpiration); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["expiration"] = transformedExpiration
	}

	transformedRequireMatchingFile, err := expandAppEngineVersionRequireMatchingFile(original["require_matching_file"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedRequireMatchingFile); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["requireMatchingFile"] = transformedRequireMatchingFile
	}

	transformedApplicationReadable, err := expandAppEngineVersionApplicationReadable(original["application_readable"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedApplicationReadable); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["applicationReadable"] = transformedApplicationReadable
	}

	return transformed, nil
}

func expandAppEngineVersionPath(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionUploadPathRegex(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionHttpHeaders(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionMimeType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionExpiration(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionRequireMatchingFile(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionApplicationReadable(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionScriptHandler(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedScriptPath, err := expandAppEngineVersionScriptPath(original["script_path"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedScriptPath); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["scriptPath"] = transformedScriptPath
	}

	return transformed, nil
}

func expandAppEngineVersionHandlers(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedUrlRegex, err := expandAppEngineVersionUrlRegex(original["url_regex"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedUrlRegex); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["urlRegex"] = transformedUrlRegex
		}

		transformedSecurityLevel, err := expandAppEngineVersionSecurityLevel(original["security_level"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSecurityLevel); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["securityLevel"] = transformedSecurityLevel
		}

		transformedLogin, err := expandAppEngineVersionLogin(original["login"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedLogin); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["login"] = transformedLogin
		}

		transformedAuthFailAction, err := expandAppEngineVersionAuthFailAction(original["auth_fail_action"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedAuthFailAction); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["authFailAction"] = transformedAuthFailAction
		}

		transformedRedirectHttpResponseCode, err := expandAppEngineVersionRedirectHttpResponseCode(original["redirect_http_response_code"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedRedirectHttpResponseCode); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["redirectHttpResponseCode"] = transformedRedirectHttpResponseCode
		}

		transformedScript, err := expandAppEngineVersionScriptHandler(original["script"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedScript); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["script"] = transformedScript
		}

		transformedStaticFiles, err := expandAppEngineVersionStaticFilesHandler(original["static_files"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedStaticFiles); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["staticFiles"] = transformedStaticFiles
		}

		req = append(req, transformed)
	}

	return req, nil
}

func expandAppEngineVersionScriptPath(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionStaticFile(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionMimeTypeFile(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionLibraries(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedName, err := expandAppEngineVersionName(original["name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["name"] = transformedName
		}

		transformedVersion, err := expandAppEngineVersionVersion(original["version"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedVersion); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["version"] = transformedVersion
		}
		req = append(req, transformed)
	}

	return req, nil
}

func expandAppEngineVersionVersion(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineUrl(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionEntrypoint(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedShell, err := expandAppEngineVersionShell(original["shell"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedShell); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["shell"] = transformedShell
	}

	return transformed, nil
}

func expandAppEngineVersionShell(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionDeployment(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedFiles, err := expandAppEngineVersionFiles(original["files"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedFiles); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["files"] = transformedFiles
	}

	transformedZip, err := expandAppEngineVersionZip(original["zip"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedZip); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["zip"] = transformedZip
	}

	return transformed, nil
}

func expandAppEngineVersionZip(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedSourceUrlProp, err := expandAppEngineVersionSourceUrl(original["source_url"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSourceUrlProp); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["sourceUrl"] = transformedSourceUrlProp
	}

	filesCountProp, err := expandAppEngineVersionFilesCount(original["files_count"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(filesCountProp); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["filesCount"] = filesCountProp
	}

	return transformed, nil
}

func expandAppEngineVersionSourceUrl(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionFilesCount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineVersionFiles(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedSourceUrlProp, err := expandAppEngineSourceUrl(original["source_url"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSourceUrlProp); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["sourceUrl"] = transformedSourceUrlProp
		}

		transformedSha1SumProp, err := expandAppEngineSha1Sum(original["sha1_sum"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSha1SumProp); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["sha1Sum"] = transformedSha1SumProp
		}

		transformedNameProp, err := expandAppEngineName(original["name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedNameProp); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["name"] = transformedNameProp
		}

		req = append(req, transformed)
	}

	return req, nil
}

func expandAppEngineSourceUrl(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineSha1Sum(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
