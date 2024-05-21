package google

import (
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const DataprocClusterAssetType string = "dataproc.googleapis.com/Cluster"

func resourceConverterDataprocCluster() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: DataprocClusterAssetType,
		Convert:   GetDataprocClusterCaiObject,
	}
}

func GetDataprocClusterCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//compute.googleapis.com/projects/{{projectId}}/regions/{{region}}/clusters/{{cluster_id]}}")
	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetDataprocClusterApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: DataprocClusterAssetType,
			Resource: &cai.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://dataproc.googleapis.com/$discovery/rest?version=v1",
				DiscoveryName:        "Cluster",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}

func GetDataprocClusterApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	projectIdProp, err := expandDataprocClusterProjectId(d.Get("project"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("project"); !tpgresource.IsEmptyValue(reflect.ValueOf(projectIdProp)) && (ok || !reflect.DeepEqual(v, projectIdProp)) {
		obj["projectId"] = projectIdProp
	}

	clusterNameProp, err := expandDataprocClusterName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(clusterNameProp)) && (ok || !reflect.DeepEqual(v, clusterNameProp)) {
		obj["clusterName"] = clusterNameProp
	}

	configProp, err := expandDataprocClusterConfig(d.Get("config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("config"); !tpgresource.IsEmptyValue(reflect.ValueOf(configProp)) && (ok || !reflect.DeepEqual(v, configProp)) {
		obj["config"] = configProp
	}

	virtualClusterConfigProp, err := expandDataprocVirtualClusterConfig(d.Get("virtual_cluster_config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("virtual_cluster_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(virtualClusterConfigProp)) && (ok || !reflect.DeepEqual(v, virtualClusterConfigProp)) {
		obj["virtualClusterConfig"] = virtualClusterConfigProp
	}

	labelsProp, err := expandDataprocClusterLabels(d.Get("effective_labels"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("effective_labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	return obj, nil
}

func expandDataprocClusterProjectId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocVirtualClusterConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataprocClusterLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}
