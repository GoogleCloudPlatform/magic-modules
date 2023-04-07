package google

import (
	"reflect"
)

func resourceConverterBigtableInstance() ResourceConverter {
	return ResourceConverter{
		AssetType: "bigtableadmin.googleapis.com/Instance",
		Convert:   GetBigtableInstanceCaiObject,
	}
}

func GetBigtableInstanceCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//bigtable.googleapis.com/projects/{{project}}/instances/{{name}}")

	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetBigtableInstanceApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: "bigtableadmin.googleapis.com/Instance",
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://bigtableadmin.googleapis.com/$discovery/rest",
				DiscoveryName:        "Instance",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetBigtableInstanceApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	nameProp, err := expandBigtableInstanceName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}

	displayNameProp, err := expandBigtableDisplayName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["name"] = nameProp
	}

	labelsProp, err := expandBigtableDisplayName(d.Get("labels"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("labels"); !isEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	return obj, nil
}

func expandBigtableInstanceName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return ReplaceVars(d, config, "projects/{{project}}/instances/{{name}}")
}

func expandBigtableDisplayName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigtableInstanceLabels(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}
