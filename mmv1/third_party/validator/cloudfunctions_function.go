package google

import (
	"reflect"
)

func GetCloudFunctionsFunctionCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//cloudfunctions.googleapis.com/projects/{{.Provider.project}}/locations/us-central1/functions/{{name}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetCloudFunctionsFunctionApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: "cloudfunctions.googleapis.com/CloudFunction",
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://cloudfunctions.googleapis.com/$discovery/rest",
				DiscoveryName:        "CloudFunction",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetCloudFunctionsFunctionApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	runtimeProp, err := expandCloudFunctionsFunctionRuntime(d.Get("runtime"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("runtime"); !isEmptyValue(reflect.ValueOf(runtimeProp)) && (ok || !reflect.DeepEqual(v, runtimeProp)) {
		obj["runtime"] = runtimeProp
	}

	nameProp, err := expandCloudFunctionsFunctionName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}

	descriptionProp, err := expandCloudFunctionsFunctionDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}

	maxInstancesProp, err := expandCloudFunctionsFunctionMaxInstances(d.Get("max_instances"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("max_instances"); !isEmptyValue(reflect.ValueOf(maxInstancesProp)) && (ok || !reflect.DeepEqual(v, maxInstancesProp)) {
		obj["maxInstances"] = maxInstancesProp
	}
	regionProp, err := expandCloudFunctionsFunctionRegion(d.Get("region"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("region"); !isEmptyValue(reflect.ValueOf(regionProp)) && (ok || !reflect.DeepEqual(v, regionProp)) {
		obj["region"] = regionProp
	}

	entryPointProp, err := expandCloudFunctionsFunctionEntryPoint(d.Get("entry_point"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("entry_point"); !isEmptyValue(reflect.ValueOf(entryPointProp)) && (ok || !reflect.DeepEqual(v, entryPointProp)) {
		obj["entryPoint"] = entryPointProp
	}

	labelsProp, err := expandCloudFunctionsFunctionLabels(d.Get("labels"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("labels"); !isEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	environmentVariablesProp, err := expandCloudFunctionsFunctionEnvironmentVariables(d.Get("environment_variables"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("environment_variables"); !isEmptyValue(reflect.ValueOf(environmentVariablesProp)) && (ok || !reflect.DeepEqual(v, environmentVariablesProp)) {
		obj["environmentVariables"] = environmentVariablesProp
	}

	buildEnvironmentVariablesProp, err := expandCloudFunctionsFunctionBuildEnvironmentVariables(d.Get("build_environment_variables"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("build_environment_variables"); !isEmptyValue(reflect.ValueOf(buildEnvironmentVariablesProp)) && (ok || !reflect.DeepEqual(v, buildEnvironmentVariablesProp)) {
		obj["buildEnvironmentVariablesProps"] = buildEnvironmentVariablesProp
	}

	availableMemoryMbProp, err := expandCloudFunctionsFunctionMemoryMb(d.Get("available_memory_mb"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("available_memory_mb"); !isEmptyValue(reflect.ValueOf(availableMemoryMbProp)) && (ok || !reflect.DeepEqual(v, availableMemoryMbProp)) {
		obj["availableMemoryMb"] = availableMemoryMbProp
	}
	
	vpcConnectorProp, err := expandCloudFunctionsFunctionVpcConnector(d.Get("vpc_connector"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("vpc_connector"); !isEmptyValue(reflect.ValueOf(vpcConnectorProp)) && (ok || !reflect.DeepEqual(v, vpcConnectorProp)) {
		obj["vpcConnector"] = vpcConnectorProp
	}

	vpcConnectorEgressSettingsProp, err := expandCloudFunctionsFunctionVpcConnectorEgressSettings(d.Get("vpc_connector_egress_settings"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("vpc_connector_egress_settings"); !isEmptyValue(reflect.ValueOf(vpcConnectorEgressSettingsProp)) && (ok || !reflect.DeepEqual(v, vpcConnectorEgressSettingsProp)) {
		obj["vpcConnectorEgressSettings"] = vpcConnectorEgressSettingsProp
	}

	ingressSettingsProp, err := expandCloudFunctionsFunctionIngressSettings(d.Get("ingress_settings"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("ingress_settings"); !isEmptyValue(reflect.ValueOf(ingressSettingsProp)) && (ok || !reflect.DeepEqual(v, ingressSettingsProp)) {
		obj["ingressSettings"] = ingressSettingsProp
	}

	serviceAccountEmailProp , err := expandCloudFunctionsFunctionServiceAccountEmail(d.Get("service_account_email"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("service_account_email"); !isEmptyValue(reflect.ValueOf(serviceAccountEmailProp)) && (ok || !reflect.DeepEqual(v, serviceAccountEmailProp)) {
		obj["serviceAccountEmail"] = serviceAccountEmailProp
	}
	return obj, nil
}

func expandCloudFunctionsFunctionRuntime(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsFunctionName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsFunctionDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsFunctionMaxInstances(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsFunctionRegion(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsFunctionEntryPoint(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsFunctionMemoryMb(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsFunctionLabels(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandCloudFunctionsFunctionEnvironmentVariables(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandCloudFunctionsFunctionBuildEnvironmentVariables(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandCloudFunctionsFunctionVpcConnector(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsFunctionVpcConnectorEgressSettings(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsFunctionIngressSettings(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsFunctionServiceAccountEmail(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
