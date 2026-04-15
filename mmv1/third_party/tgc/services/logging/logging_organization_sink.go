package logging

import (
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

func ResourceConverterLogOrganizationSink() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: logSinkAssetType,
		Convert:   GetLogOrganizationSinkCaiObject,
	}
}

func GetLogOrganizationSinkCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//logging.googleapis.com/organizations/{{org_id}}/sinks/{{name}}")
	if err != nil {
		return []cai.Asset{}, err
	}
	obj, err := GetLogOrganizationSinkApiObject(d, config)
	if err != nil {
		return []cai.Asset{}, err
	}
	return []cai.Asset{{
		Name: name,
		Type: logSinkAssetType,
		Resource: &cai.AssetResource{
			Version:              "v2",
			DiscoveryDocumentURI: "https://logging.googleapis.com/$discovery/rest",
			DiscoveryName:        "LogSink",
			Data:                 obj,
		},
	}}, nil
}

func GetLogOrganizationSinkApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	nameProp, err := expandLogOrganizationSinkName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}

	destinationProp, err := expandLogOrganizationSinkDestination(d.Get("destination"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("destination"); !tpgresource.IsEmptyValue(reflect.ValueOf(destinationProp)) && (ok || !reflect.DeepEqual(v, destinationProp)) {
		obj["destination"] = destinationProp
	}

	filterProp, err := expandLogOrganizationSinkFilter(d.Get("filter"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("filter"); !tpgresource.IsEmptyValue(reflect.ValueOf(filterProp)) && (ok || !reflect.DeepEqual(v, filterProp)) {
		obj["filter"] = filterProp
	}

	descriptionProp, err := expandLogOrganizationSinkDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}

	disabledProp, err := expandLogOrganizationSinkDisabled(d.Get("disabled"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("disabled"); !tpgresource.IsEmptyValue(reflect.ValueOf(disabledProp)) && (ok || !reflect.DeepEqual(v, disabledProp)) {
		obj["disabled"] = disabledProp
	}

	exclusionsProp, err := expandLogOrganizationSinkExclusions(d.Get("exclusions"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("exclusions"); !tpgresource.IsEmptyValue(reflect.ValueOf(exclusionsProp)) && (ok || !reflect.DeepEqual(v, exclusionsProp)) {
		obj["exclusions"] = exclusionsProp
	}

	includeChildrenProp, err := expandLogOrganizationSinkIncludeChildren(d.Get("include_children"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("include_children"); !tpgresource.IsEmptyValue(reflect.ValueOf(includeChildrenProp)) && (ok || !reflect.DeepEqual(v, includeChildrenProp)) {
		obj["includeChildren"] = includeChildrenProp
	}

	interceptChildrenProp, err := expandLogOrganizationSinkInterceptChildren(d.Get("intercept_children"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("intercept_children"); !tpgresource.IsEmptyValue(reflect.ValueOf(interceptChildrenProp)) && (ok || !reflect.DeepEqual(v, interceptChildrenProp)) {
		obj["interceptChildren"] = interceptChildrenProp
	}

	bigqueryOptionsProp, err := expandLogOrganizationSinkBigqueryOptions(d.Get("bigquery_options"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("bigquery_options"); !tpgresource.IsEmptyValue(reflect.ValueOf(bigqueryOptionsProp)) && (ok || !reflect.DeepEqual(v, bigqueryOptionsProp)) {
		obj["bigqueryOptions"] = bigqueryOptionsProp
	}

	return obj, nil
}

func expandLogOrganizationSinkName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogOrganizationSinkDestination(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogOrganizationSinkFilter(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogOrganizationSinkDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogOrganizationSinkDisabled(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogOrganizationSinkIncludeChildren(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogOrganizationSinkInterceptChildren(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogOrganizationSinkExclusions(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
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

		transformedName, err := expandLogOrganizationSinkExclusionsName(original["name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["name"] = transformedName
		}

		transformedDescription, err := expandLogOrganizationSinkExclusionsDescription(original["description"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedDescription); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["description"] = transformedDescription
		}

		transformedFilter, err := expandLogOrganizationSinkExclusionsFilter(original["filter"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedFilter); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["filter"] = transformedFilter
		}

		transformedDisabled, err := expandLogOrganizationSinkExclusionsDisabled(original["disabled"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedDisabled); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["disabled"] = transformedDisabled
		}

		req = append(req, transformed)
	}

	return req, nil
}

func expandLogOrganizationSinkExclusionsName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogOrganizationSinkExclusionsDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogOrganizationSinkExclusionsFilter(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogOrganizationSinkExclusionsDisabled(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogOrganizationSinkBigqueryOptions(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedUsePartitionedTables, err := expandLogOrganizationSinkBigqueryOptionsUsePartitionedTables(original["use_partitioned_tables"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedUsePartitionedTables); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["usePartitionedTables"] = transformedUsePartitionedTables
	}

	return transformed, nil
}

func expandLogOrganizationSinkBigqueryOptionsUsePartitionedTables(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
