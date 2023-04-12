package google

import (
	"reflect"
)

func resourceConverterBigtableCluster() ResourceConverter {
	return ResourceConverter{
		AssetType: "bigtableadmin.googleapis.com/Cluster",
		Convert:   GetBigtableClusterCaiObject,
	}
}

func GetBigtableClusterCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {

	objs, err := GetBigtableClusterApiObjects(d, config)

	if err != nil {
		return []Asset{}, err
	}

	assets := []Asset{}
	for _, obj := range objs {
		name, err := assetName(d, config, "//bigtable.googleapis.com/projects/{{project}}/instances/{{name}}/clusters/{{cluster_id}}")
		if err != nil {
			return []Asset{}, err
		}

		asset := Asset{
			Name: name,
			Type: "bigtableadmin.googleapis.com/Cluster",
			Resource: &AssetResource{
				Version:              "v2",
				DiscoveryDocumentURI: "https://bigtableadmin.googleapis.com/$discovery/rest",
				DiscoveryName:        "Cluster",
				Data:                 obj,
			},
		}
		assets = append(assets, asset)
	}
	return assets, nil
}

func GetBigtableClusterApiObjects(d TerraformResourceData, config *Config) ([]map[string]interface{}, error) {
	return expandBigtableClusters(d.Get("cluster"), d, config)

}

func expandBigtableClusters(v interface{}, d TerraformResourceData, config *Config) ([]map[string]interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}

	transformedEntries := []map[string]interface{}{}

	for _, raw := range l {
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedLocation, err := expandBigtableClusterLocation(original["zone"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedLocation); val.IsValid() && !isEmptyValue(val) {
			transformed["location"] = transformedLocation
		}

		transformedServerNodes, err := expandBigtableClusterServerNodes(original["num_nodes"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedServerNodes); val.IsValid() && !isEmptyValue(val) {
			transformed["serverNodes"] = transformedServerNodes
		}

		transformedStorageType, err := expandBigtableClusterDefaultStorageType(original["storage_type"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedStorageType); val.IsValid() && !isEmptyValue(val) {
			transformed["defaultStorageType"] = transformedStorageType
		}

		transformedName, err := expandBigtableClusterName(original["cluster_id"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedName); val.IsValid() && !isEmptyValue(val) {
			transformed["name"] = transformedName
		}
		transformedEntries = append(transformedEntries, transformed)
	}

	return transformedEntries, nil
}

func expandBigtableClusterLocation(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigtableClusterServerNodes(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigtableClusterDefaultStorageType(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigtableClusterName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	cluster, err := ReplaceVars(d, config, "projects/{{project}}/instances/{{name}}/clusters/")
	if err != nil {
		return nil, err
	}
	return cluster + v.(string), nil
}
