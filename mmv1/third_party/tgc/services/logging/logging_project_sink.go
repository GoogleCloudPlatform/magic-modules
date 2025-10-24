package logging

import (
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const logSinkAssetType string = "logging.googleapis.com/LogSink"

func ResourceConverterLogProjectSink() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: logSinkAssetType,
		Convert:   GetLogProjectSinkCaiObject,
	}
}

func GetLogProjectSinkCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//logging.googleapis.com/projects/{{project}}/sinks/{{name}}")
	if err != nil {
		return []cai.Asset{}, err
	}
	obj, err := GetLogProjectSinkApiObject(d, config)
	if err != nil {
		return []cai.Asset{}, err
	}
	return []cai.Asset{{
		Name: name,
		Type: logSinkAssetType,
		Resource: &cai.AssetResource{
			Version:              "v2",
			DiscoveryDocumentURI: "https://logging.googleapis.com/$discovery/rest?version=v2",
			DiscoveryName:        "LogSink",
			Data:                 obj,
		},
	}}, nil
}

func GetLogProjectSinkApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	nameProp, err := expandLogProjectSinkName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}

	destinationProp, err := expandLogProjectSinkDestination(d.Get("destination"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("destination"); !tpgresource.IsEmptyValue(reflect.ValueOf(destinationProp)) && (ok || !reflect.DeepEqual(v, destinationProp)) {
		obj["destination"] = destinationProp
	}

	filterProp, err := expandLogProjectSinkFilter(d.Get("filter"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("filter"); !tpgresource.IsEmptyValue(reflect.ValueOf(filterProp)) && (ok || !reflect.DeepEqual(v, filterProp)) {
		obj["filter"] = filterProp
	}

	descriptionProp, err := expandLogProjectSinkDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}

	disabledProp, err := expandLogProjectSinkDisabled(d.Get("disabled"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("disabled"); !tpgresource.IsEmptyValue(reflect.ValueOf(disabledProp)) && (ok || !reflect.DeepEqual(v, disabledProp)) {
		obj["disabled"] = disabledProp
	}

	exclusionsProp, err := expandLogProjectSinkExclusions(d.Get("exclusions"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("exclusions"); !tpgresource.IsEmptyValue(reflect.ValueOf(exclusionsProp)) && (ok || !reflect.DeepEqual(v, exclusionsProp)) {
		obj["exclusions"] = exclusionsProp
	}

	bigqueryOptionsProp, err := expandLogProjectSinkBigqueryOptions(d.Get("bigquery_options"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("bigquery_options"); !tpgresource.IsEmptyValue(reflect.ValueOf(bigqueryOptionsProp)) && (ok || !reflect.DeepEqual(v, bigqueryOptionsProp)) {
		obj["bigqueryOptions"] = bigqueryOptionsProp
	}

	return obj, nil
}

func expandLogProjectSinkName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogProjectSinkDestination(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogProjectSinkFilter(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogProjectSinkDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogProjectSinkDisabled(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogProjectSinkExclusions(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l, ok := v.([]interface{})
	if !ok {
		return nil, nil
	}
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedName, err := expandLogProjectSinkExclusionsName(original["name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["name"] = transformedName
		}

		transformedDescription, err := expandLogProjectSinkExclusionsDescription(original["description"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedDescription); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["description"] = transformedDescription
		}

		transformedFilter, err := expandLogProjectSinkExclusionsFilter(original["filter"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedFilter); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["filter"] = transformedFilter
		}

		transformedDisabled, err := expandLogProjectSinkExclusionsDisabled(original["disabled"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedDisabled); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["disabled"] = transformedDisabled
		}

		req = append(req, transformed)
	}

	return req, nil
}

func expandLogProjectSinkExclusionsName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogProjectSinkExclusionsDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogProjectSinkExclusionsFilter(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogProjectSinkExclusionsDisabled(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogProjectSinkBigqueryOptions(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedUsePartitionedTables, err := expandLogProjectSinkBigqueryOptionsUsePartitionedTables(original["use_partitioned_tables"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedUsePartitionedTables); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["usePartitionedTables"] = transformedUsePartitionedTables
	}

	return transformed, nil
}

func expandLogProjectSinkBigqueryOptionsUsePartitionedTables(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
