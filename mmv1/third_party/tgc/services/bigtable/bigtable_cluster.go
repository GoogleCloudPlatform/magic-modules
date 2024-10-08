package bigtable

import (
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

func ResourceConverterBigtableCluster() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: "bigtableadmin.googleapis.com/Cluster",
		Convert:   GetBigtableClusterCaiObject,
	}
}

func GetBigtableClusterCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {

	objs, err := GetBigtableClusterApiObjects(d, config)

	if err != nil {
		return []cai.Asset{}, err
	}

	assets := []cai.Asset{}
	for _, obj := range objs {
		name, err := cai.AssetName(d, config, "//bigtable.googleapis.com/projects/{{project}}/instances/{{name}}/clusters/{{cluster_id}}")
		if err != nil {
			return []cai.Asset{}, err
		}

		asset := cai.Asset{
			Name: name,
			Type: "bigtableadmin.googleapis.com/Cluster",
			Resource: &cai.AssetResource{
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

func GetBigtableClusterApiObjects(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]map[string]interface{}, error) {
	return expandBigtableClusters(d.Get("cluster"), d, config)

}

func expandBigtableClusters(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]map[string]interface{}, error) {
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
		} else if val := reflect.ValueOf(transformedLocation); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["location"] = transformedLocation
		}

		transformedServerNodes, err := expandBigtableClusterServerNodes(original["num_nodes"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedServerNodes); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["serverNodes"] = transformedServerNodes
		}

		transformedStorageType, err := expandBigtableClusterDefaultStorageType(original["storage_type"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedStorageType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["defaultStorageType"] = transformedStorageType
		}

		transformedName, err := expandBigtableClusterName(original["cluster_id"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["name"] = transformedName
		}
		transformedEntries = append(transformedEntries, transformed)
	}

	return transformedEntries, nil
}

func expandBigtableClusterLocation(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandBigtableClusterServerNodes(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandBigtableClusterDefaultStorageType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandBigtableClusterName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	cluster, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/instances/{{name}}/clusters/")
	if err != nil {
		return nil, err
	}
	return cluster + v.(string), nil
}
