package logging

import (
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const LogSinkAssetType string = "logging.googleapis.com/LogSink"

func ResourceConverterLogFolderSink() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: LogSinkAssetType,
		Convert:   GetLogFolderSinkCaiObject,
	}
}

func GetLogFolderSinkCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//logging.googleapis.com/folders/{{folder}}/sinks/{{name}}")
	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetLogFolderSinkApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: LogSinkAssetType,
			Resource: &cai.AssetResource{
				Version:              "v2",
				DiscoveryDocumentURI: "https://logging.googleapis.com/$discovery/rest",
				DiscoveryName:        "LogSink",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}

func GetLogFolderSinkApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	nameProp, err := expandLogFolderSinkName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}

	destinationProp, err := expandLogFolderSinkDestination(d.Get("destination"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("destination"); !tpgresource.IsEmptyValue(reflect.ValueOf(destinationProp)) && (ok || !reflect.DeepEqual(v, destinationProp)) {
		obj["destination"] = destinationProp
	}

	filterProp, err := expandLogFolderSinkFilter(d.Get("filter"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("filter"); !tpgresource.IsEmptyValue(reflect.ValueOf(filterProp)) && (ok || !reflect.DeepEqual(v, filterProp)) {
		obj["filter"] = filterProp
	}

	descriptionProp, err := expandLogFolderSinkDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}

	disabledProp, err := expandLogFolderSinkDisabled(d.Get("disabled"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("disabled"); !tpgresource.IsEmptyValue(reflect.ValueOf(disabledProp)) && (ok || !reflect.DeepEqual(v, disabledProp)) {
		obj["disabled"] = disabledProp
	}

	exclusionsProp, err := expandLogFolderSinkExclusions(d.Get("exclusions"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("exclusions"); !tpgresource.IsEmptyValue(reflect.ValueOf(exclusionsProp)) && (ok || !reflect.DeepEqual(v, exclusionsProp)) {
		obj["exclusions"] = exclusionsProp
	}

	includeChildrenProp, err := expandLogFolderSinkIncludeChildren(d.Get("include_children"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("include_children"); !tpgresource.IsEmptyValue(reflect.ValueOf(includeChildrenProp)) && (ok || !reflect.DeepEqual(v, includeChildrenProp)) {
		obj["includeChildren"] = includeChildrenProp
	}

	bigqueryOptionsProp, err := expandLogFolderSinkBigqueryOptions(d.Get("bigquery_options"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("bigquery_options"); !tpgresource.IsEmptyValue(reflect.ValueOf(bigqueryOptionsProp)) && (ok || !reflect.DeepEqual(v, bigqueryOptionsProp)) {
		obj["bigqueryOptions"] = bigqueryOptionsProp
	}

	return obj, nil
}

func expandLogFolderSinkName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogFolderSinkDestination(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogFolderSinkFilter(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogFolderSinkDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogFolderSinkDisabled(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogFolderSinkExclusions(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
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

		transformedName, err := expandLogFolderSinkExclusionsName(original["name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["name"] = transformedName
		}

		transformedDescription, err := expandLogFolderSinkExclusionsDescription(original["description"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedDescription); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["description"] = transformedDescription
		}

		transformedFilter, err := expandLogFolderSinkExclusionsFilter(original["filter"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedFilter); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["filter"] = transformedFilter
		}

		transformedDisabled, err := expandLogFolderSinkExclusionsDisabled(original["disabled"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedDisabled); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["disabled"] = transformedDisabled
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandLogFolderSinkExclusionsName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogFolderSinkExclusionsDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogFolderSinkExclusionsFilter(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogFolderSinkExclusionsDisabled(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogFolderSinkIncludeChildren(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogFolderSinkBigqueryOptions(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedUsePartitionedTables, err := expandLogFolderSinkBigqueryOptionsUsePartitionedTables(original["use_partitioned_tables"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedUsePartitionedTables); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["usePartitionedTables"] = transformedUsePartitionedTables
	}

	return transformed, nil
}

func expandLogFolderSinkBigqueryOptionsUsePartitionedTables(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
