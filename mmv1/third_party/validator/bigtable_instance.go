package google

import (
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
)

func resourceConverterBigtableInstance() ResourceConverter {
	return ResourceConverter{
		AssetType: "bigtableadmin.googleapis.com/Instance",
		Convert:   GetBigtableInstanceCaiObject,
	}
}

func GetBigtableInstanceCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]Asset, error) {
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

func GetBigtableInstanceApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	nameProp, err := expandBigtableInstanceName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}

	displayNameProp, err := expandBigtableDisplayName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["name"] = nameProp
	}

	labelsProp, err := expandBigtableDisplayName(d.Get("labels"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	return obj, nil
}

func expandBigtableInstanceName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return tpgresource.ReplaceVars(d, config, "projects/{{project}}/instances/{{name}}")
}

func expandBigtableDisplayName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandBigtableInstanceLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}
