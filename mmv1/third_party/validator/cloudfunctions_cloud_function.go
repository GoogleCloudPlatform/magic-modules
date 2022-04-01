package google

import "reflect"

const CloudFunctionsCloudFunctionAssetType string = "cloudfunctions.googleapis.com/CloudFunction"

func resourceConverterCloudFunctionsCloudFunction() ResourceConverter {
	return ResourceConverter{
		AssetType: CloudFunctionsCloudFunctionAssetType,
		Convert:   GetCloudFunctionsCloudFunctionCaiObject,
	}
}

func GetCloudFunctionsCloudFunctionCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//cloudfunctions.googleapis.com/projects/{{project}}/locations/{{region}}/functions/{{name}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetCloudFunctionsCloudFunctionApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: CloudFunctionsCloudFunctionAssetType,
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/cloudfunctions/v1/rest",
				DiscoveryName:        "CloudFunction",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetCloudFunctionsCloudFunctionApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	nameProp, err := expandCloudFunctionsCloudFunctionName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	descriptionProp, err := expandCloudFunctionsCloudFunctionDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	entryPointProp, err := expandCloudFunctionsCloudFunctionEntryPoint(d.Get("entry_point"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("entry_point"); !isEmptyValue(reflect.ValueOf(entryPointProp)) && (ok || !reflect.DeepEqual(v, entryPointProp)) {
		obj["entryPoint"] = entryPointProp
	}
	runtimeProp, err := expandCloudFunctionsCloudFunctionRuntime(d.Get("runtime"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("runtime"); !isEmptyValue(reflect.ValueOf(runtimeProp)) && (ok || !reflect.DeepEqual(v, runtimeProp)) {
		obj["runtime"] = runtimeProp
	}
	timeoutProp, err := expandCloudFunctionsCloudFunctionTimeout(d.Get("timeout"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("timeout"); !isEmptyValue(reflect.ValueOf(timeoutProp)) && (ok || !reflect.DeepEqual(v, timeoutProp)) {
		obj["timeout"] = timeoutProp
	}
	availableMemoryMbProp, err := expandCloudFunctionsCloudFunctionAvailableMemoryMb(d.Get("available_memory_mb"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("available_memory_mb"); !isEmptyValue(reflect.ValueOf(availableMemoryMbProp)) && (ok || !reflect.DeepEqual(v, availableMemoryMbProp)) {
		obj["availableMemoryMb"] = availableMemoryMbProp
	}
	labelsProp, err := expandCloudFunctionsCloudFunctionLabels(d.Get("labels"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("labels"); !isEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}
	environmentVariablesProp, err := expandCloudFunctionsCloudFunctionEnvironmentVariables(d.Get("environment_variables"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("environment_variables"); !isEmptyValue(reflect.ValueOf(environmentVariablesProp)) && (ok || !reflect.DeepEqual(v, environmentVariablesProp)) {
		obj["environmentVariables"] = environmentVariablesProp
	}
	sourceArchiveUrlProp, err := expandCloudFunctionsCloudFunctionSourceArchiveUrl(d.Get("source_archive_url"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("source_archive_url"); !isEmptyValue(reflect.ValueOf(sourceArchiveUrlProp)) && (ok || !reflect.DeepEqual(v, sourceArchiveUrlProp)) {
		obj["sourceArchiveUrl"] = sourceArchiveUrlProp
	}
	sourceUploadUrlProp, err := expandCloudFunctionsCloudFunctionSourceUploadUrl(d.Get("source_upload_url"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("source_upload_url"); !isEmptyValue(reflect.ValueOf(sourceUploadUrlProp)) && (ok || !reflect.DeepEqual(v, sourceUploadUrlProp)) {
		obj["sourceUploadUrl"] = sourceUploadUrlProp
	}
	sourceRepositoryProp, err := expandCloudFunctionsCloudFunctionSourceRepository(d.Get("source_repository"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("source_repository"); !isEmptyValue(reflect.ValueOf(sourceRepositoryProp)) && (ok || !reflect.DeepEqual(v, sourceRepositoryProp)) {
		obj["sourceRepository"] = sourceRepositoryProp
	}
	httpsTriggerProp, err := expandCloudFunctionsCloudFunctionHttpsTriggerUrl(d.Get("https_trigger_url"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("https_trigger_url"); !isEmptyValue(reflect.ValueOf(httpsTriggerProp)) && (ok || !reflect.DeepEqual(v, httpsTriggerProp)) {
		obj["httpsTriggerUrl"] = httpsTriggerProp
	}
	eventTriggerProp, err := expandCloudFunctionsCloudFunctionEventTrigger(d.Get("event_trigger"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("event_trigger"); !isEmptyValue(reflect.ValueOf(eventTriggerProp)) && (ok || !reflect.DeepEqual(v, eventTriggerProp)) {
		obj["eventTrigger"] = eventTriggerProp
	}
	locationProp, err := expandCloudFunctionsCloudFunctionRegion(d.Get("region"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("region"); !isEmptyValue(reflect.ValueOf(locationProp)) && (ok || !reflect.DeepEqual(v, locationProp)) {
		obj["location"] = locationProp
	}
	trigger_httpProp, err := expandCloudFunctionsCloudFunctionTriggerHttp(d.Get("trigger_http"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("trigger_http"); !isEmptyValue(reflect.ValueOf(trigger_httpProp)) && (ok || !reflect.DeepEqual(v, trigger_httpProp)) {
		obj["trigger_http"] = trigger_httpProp
	}

	return obj, nil
}

func expandCloudFunctionsCloudFunctionName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsCloudFunctionDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsCloudFunctionEntryPoint(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsCloudFunctionRuntime(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsCloudFunctionTimeout(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsCloudFunctionAvailableMemoryMb(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsCloudFunctionLabels(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandCloudFunctionsCloudFunctionEnvironmentVariables(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandCloudFunctionsCloudFunctionSourceArchiveUrl(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsCloudFunctionSourceUploadUrl(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsCloudFunctionSourceRepository(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedUrl, err := expandCloudFunctionsCloudFunctionSourceRepositoryUrl(original["url"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedUrl); val.IsValid() && !isEmptyValue(val) {
		transformed["url"] = transformedUrl
	}

	transformedDeployedUrl, err := expandCloudFunctionsCloudFunctionSourceRepositoryDeployedUrl(original["deployed_url"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDeployedUrl); val.IsValid() && !isEmptyValue(val) {
		transformed["deployedUrl"] = transformedDeployedUrl
	}

	return transformed, nil
}

func expandCloudFunctionsCloudFunctionSourceRepositoryUrl(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsCloudFunctionSourceRepositoryDeployedUrl(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsCloudFunctionHttpsTriggerUrl(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsCloudFunctionEventTrigger(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedEventType, err := expandCloudFunctionsCloudFunctionEventTriggerEventType(original["event_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEventType); val.IsValid() && !isEmptyValue(val) {
		transformed["eventType"] = transformedEventType
	}

	transformedResource, err := expandCloudFunctionsCloudFunctionEventTriggerResource(original["resource"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedResource); val.IsValid() && !isEmptyValue(val) {
		transformed["resource"] = transformedResource
	}

	transformedService, err := expandCloudFunctionsCloudFunctionEventTriggerService(original["service"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedService); val.IsValid() && !isEmptyValue(val) {
		transformed["service"] = transformedService
	}

	return transformed, nil
}

func expandCloudFunctionsCloudFunctionEventTriggerEventType(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsCloudFunctionEventTriggerResource(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsCloudFunctionEventTriggerService(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsCloudFunctionRegion(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudFunctionsCloudFunctionTriggerHttp(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
