package logging

import (
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const logOrganizationBucketAssetType string = "logging.googleapis.com/LogBucket"

func ResourceConverterLogOrganizationBucket() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: logOrganizationBucketAssetType,
		Convert:   GetLogOrganizationBucketCaiObject,
	}
}

func GetLogOrganizationBucketCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//logging.googleapis.com/projects/{{project}}/locations/{{location}}/buckets/{{bucket_id}}")
	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetLogOrganizationBucketApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: logOrganizationBucketAssetType,
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

func GetLogOrganizationBucketApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	organizationProp, err := expandLogOrganizationBucketOrganizationId(d.Get("organization"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("organization"); !tpgresource.IsEmptyValue(reflect.ValueOf(organizationProp)) && (ok || !reflect.DeepEqual(v, organizationProp)) {
		obj["id"] = organizationProp
	}

	nameProp, err := expandLogOrganizationBucketName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}

	bucketIdProp, err := expandLogOrganizationBucketBucketId(d.Get("bucket_id"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("bucket_id"); !tpgresource.IsEmptyValue(reflect.ValueOf(bucketIdProp)) && (ok || !reflect.DeepEqual(v, bucketIdProp)) {
		obj["bucketId"] = bucketIdProp
	}

	locationProp, err := expandLogOrganizationBucketLocation(d.Get("location"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("location"); !tpgresource.IsEmptyValue(reflect.ValueOf(locationProp)) && (ok || !reflect.DeepEqual(v, locationProp)) {
		obj["location"] = locationProp
	}

	descriptionProp, err := expandLogOrganizationBucketDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}

	retentionDaysProp, err := expandLogOrganizationBucketRetentionDays(d.Get("retention_days"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("retention_days"); !tpgresource.IsEmptyValue(reflect.ValueOf(retentionDaysProp)) && (ok || !reflect.DeepEqual(v, retentionDaysProp)) {
		obj["retentionDays"] = retentionDaysProp
	}

	indexConfigsProp, err := expandLogOrganizationBucketIndexConfigs(d.Get("index_configs"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("index_configs"); !tpgresource.IsEmptyValue(reflect.ValueOf(indexConfigsProp)) && (ok || !reflect.DeepEqual(v, indexConfigsProp)) {
		obj["indexConfigs"] = indexConfigsProp
	}

	lifecycleStateProp, err := expandLogOrganizationBucketLifecycleState(d.Get("lifecycle_state"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("lifecycle_state"); !tpgresource.IsEmptyValue(reflect.ValueOf(lifecycleStateProp)) && (ok || !reflect.DeepEqual(v, lifecycleStateProp)) {
		obj["lifecycleState"] = lifecycleStateProp
	}

	return obj, nil
}

func expandLogOrganizationBucketOrganizationId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	v, err := tpgresource.ReplaceVars(d, config, "organizations/{{organization}}/locations/{{location}}/buckets/{{bucket_id}}")
	if err != nil {
		return nil, err
	}

	return v, nil
}

func expandLogOrganizationBucketName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogOrganizationBucketLifecycleState(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogOrganizationBucketIndexConfigs(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedFieldPath, err := expandLogOrganizationBucketFieldPath(original["field_path"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedFieldPath); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["fieldPath"] = transformedFieldPath
		}

		transformedType, err := expandLogOrganizationBucketType(original["type"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["type"] = transformedType
		}

		req = append(req, transformed)
	}

	return req, nil
}

func expandLogOrganizationBucketType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogOrganizationBucketFieldPath(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogOrganizationBucketRetentionDays(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogOrganizationBucketDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogOrganizationBucketLocation(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLogOrganizationBucketBucketId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
