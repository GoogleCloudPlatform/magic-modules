package apikeys

import (
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const ApikeysKeyAssetType string = "apikeys.googleapis.com/Key"

func ResourceConverterApikeysKey() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: ApikeysKeyAssetType,
		Convert:   GetApikeysKeyCaiObject,
	}
}

func GetApikeysKeyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//apikeys.googleapis.com/v2/projects/{{project}}/locations/global/keys/{{key}}")
	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetApikeysKeyApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: ApikeysKeyAssetType,
			Resource: &cai.AssetResource{
				Version:              "v2",
				DiscoveryDocumentURI: "https://apikeys.googleapis.com/$discovery/rest?version=v2",
				DiscoveryName:        "Apikeyskey",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}

func GetApikeysKeyApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	uidProp, err := expandApikeysKeyUid(d.Get("uid"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("uid"); !tpgresource.IsEmptyValue(reflect.ValueOf(uidProp)) && (ok || !reflect.DeepEqual(v, uidProp)) {
		obj["uid"] = uidProp
	}

	displayNameProp, err := expandApikeysKeyDisplayName(d.Get("displayName"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("displayName"); !tpgresource.IsEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}

	keyStringProp, err := expandApikeysKeyKeyString(d.Get("keyString"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("keyString"); !tpgresource.IsEmptyValue(reflect.ValueOf(keyStringProp)) && (ok || !reflect.DeepEqual(v, keyStringProp)) {
		obj["keyString"] = keyStringProp
	}

	createTimeProp, err := expandApikeysKeyCreateTime(d.Get("createTime"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("createTime"); !tpgresource.IsEmptyValue(reflect.ValueOf(createTimeProp)) && (ok || !reflect.DeepEqual(v, createTimeProp)) {
		obj["createTime"] = createTimeProp
	}

	updateTimeProp, err := expandApikeysKeyUpdateTime(d.Get("updateTime"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("updateTime"); !tpgresource.IsEmptyValue(reflect.ValueOf(updateTimeProp)) && (ok || !reflect.DeepEqual(v, updateTimeProp)) {
		obj["updateTime"] = updateTimeProp
	}

	deleteTimeProp, err := expandApikeysKeyDeleteTime(d.Get("deleteTime"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("deleteTime"); !tpgresource.IsEmptyValue(reflect.ValueOf(deleteTimeProp)) && (ok || !reflect.DeepEqual(v, deleteTimeProp)) {
		obj["deleteTime"] = deleteTimeProp
	}

	restrictionsProp, err := expandApikeysKeyRestrictions(d.Get("restrictions"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("restrictions"); !tpgresource.IsEmptyValue(reflect.ValueOf(restrictionsProp)) && (ok || !reflect.DeepEqual(v, restrictionsProp)) {
		obj["restrictions"] = restrictionsProp
	}

	etagProp, err := expandApikeysKeyDEtag(d.Get("etag"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("etag"); !tpgresource.IsEmptyValue(reflect.ValueOf(etagProp)) && (ok || !reflect.DeepEqual(v, etagProp)) {
		obj["etag"] = etagProp
	}

	return obj, nil
}

func expandApikeysKeyUid(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApikeysKeyDisplayName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApikeysKeyKeyString(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApikeysKeyCreateTime(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApikeysKeyUpdateTime(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApikeysKeyDeleteTime(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApikeysKeyRestrictions(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {

	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedAndroidKeyRestrictions, err := expandApikeysKeyAndroidKeyRestriction(original["android_key_restrictions"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAndroidKeyRestrictions); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["android_key_restrictions"] = transformedAndroidKeyRestrictions
	}

	transformedApiTargets, err := expandApikeysKeyApiTargets(original["api_targets"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedApiTargets); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["api_targets"] = transformedApiTargets
	}

	return transformed, nil
}
func expandApikeysKeyAndroidKeyRestriction(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedAllowedServices, err := expandApikeysKeyAllowedApplications(original["allowed_applications"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAllowedServices); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["allowed_applications"] = transformedAllowedServices
	}

	return transformed, nil
}

func expandApikeysKeyAllowedApplications(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedPackageName, err := expandApikeysKeyPackageName(original["package_name"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPackageName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["package_name"] = transformedPackageName
	}

	transformedSha1Fingerprint, err := expandApikeysKeySha1Fingerprint(original["sha1_fingerprint"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSha1Fingerprint); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["sha1_fingerprint"] = transformedSha1Fingerprint
	}

	return transformed, nil
}

func expandApikeysKeyPackageName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApikeysKeySha1Fingerprint(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApikeysKeyApiTargets(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedService, err := expandApikeysKeyService(original["service"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedService); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["service"] = transformedService
	}

	transformedMethods, err := expandApikeysKeyMethods(original["methods"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMethods); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["methods"] = transformedMethods
	}

	return transformed, nil
}

func expandApikeysKeyService(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApikeysKeyMethods(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return cai.ConvertInterfaceToStringArray(v.([]interface{})), nil
}

func expandApikeysKeyDEtag(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
