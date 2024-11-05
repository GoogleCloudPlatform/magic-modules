package logging

import (
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const logFolderBucketAssetType string = "logging.googleapis.com/LogBucket"

func ResourceConverterLogFolderBucket() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: logFolderBucketAssetType,
		Convert:   GetLogFolderBucketCaiObject,
	}
}

func GetLogFolderBucketCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//logging.googleapis.com/projects/{{project}}/locations/{{location}}/buckets/{{bucket_id}}")
	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetLogFolderBucketApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: logFolderBucketAssetType,
			Resource: &cai.AssetResource{
				Version:              "v2",
				DiscoveryDocumentURI: "https://logging.googleapis.com/$discovery/rest?version=v2",
				DiscoveryName:        "LogBucket",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}

func GetLogFolderBucketApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	folderProp, err := expandLogFolderBucketFolderId(d.Get("folder"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("folder"); !tpgresource.IsEmptyValue(reflect.ValueOf(folderProp)) && (ok || !reflect.DeepEqual(v, folderProp)) {
		obj["id"] = folderProp
	}

	nameProp, err := expandLogFolderBucketName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}

	bucketIdProp, err := expandLogFolderBucketBucketId(d.Get("bucket_id"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("bucket_id"); !tpgresource.IsEmptyValue(reflect.ValueOf(bucketIdProp)) && (ok || !reflect.DeepEqual(v, bucketIdProp)) {
		obj["bucketId"] = bucketIdProp
	}

	locationProp, err := expandLogFolderBucketLocation(d.Get("location"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("location"); !tpgresource.IsEmptyValue(reflect.ValueOf(locationProp)) && (ok || !reflect.DeepEqual(v, locationProp)) {
		obj["location"] = locationProp
	}

	descriptionProp, err := expandLogFolderBucketDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}

	retentionDaysProp, err := expandLogFolderBucketRetentionDays(d.Get("retention_days"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("retention_days"); !tpgresource.IsEmptyValue(reflect.ValueOf(retentionDaysProp)) && (ok || !reflect.DeepEqual(v, retentionDaysProp)) {
		obj["retentionDays"] = retentionDaysProp
	}

	indexConfigsProp, err := expandLogFolderBucketIndexConfigs(d.Get("index_configs"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("index_configs"); !tpgresource.IsEmptyValue(reflect.ValueOf(indexConfigsProp)) && (ok || !reflect.DeepEqual(v, indexConfigsProp)) {
		obj["indexConfigs"] = indexConfigsProp
	}

	lifecycleStateProp, err := expandLogFolderBucketLifecycleState(d.Get("lifecycle_state"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("lifecycle_state"); !tpgresource.IsEmptyValue(reflect.ValueOf(lifecycleStateProp)) && (ok || !reflect.DeepEqual(v, lifecycleStateProp)) {
		obj["lifecycleState"] = lifecycleStateProp
	}

	return obj, nil
}

func expandLogFolderBucketFolderId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	v, err := tpgresource.ReplaceVars(d, config, "folders/{{folder}}/locations/{{location}}/buckets/{{bucket_id}}")
	if err != nil {
		return nil, err
	}

	return v, nil
}

func expandLogFolderBucketName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogFolderBucketLifecycleState(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogFolderBucketIndexConfigs(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedFieldPath, err := expandLogFolderBucketFieldPath(original["field_path"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedFieldPath); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["fieldPath"] = transformedFieldPath
		}

		transformedType, err := expandLogFolderBucketType(original["type"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["type"] = transformedType
		}

		req = append(req, transformed)
	}

	return req, nil
}

func expandLogFolderBucketType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogFolderBucketFieldPath(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogFolderBucketRetentionDays(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogFolderBucketDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogFolderBucketLocation(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogFolderBucketBucketId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogFolderBucketFolder(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
