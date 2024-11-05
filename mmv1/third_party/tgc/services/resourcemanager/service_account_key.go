package resourcemanager

import (
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const ServiceAccountKeyAssetType string = "iam.googleapis.com/ServiceAccountKey"

func ResourceConverterServiceAccountKey() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: ServiceAccountKeyAssetType,
		Convert:   GetServiceAccountKeyCaiObject,
	}
}

func GetServiceAccountKeyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//iam.googleapis.com/projects/{{project}}/serviceAccounts/{{account}}/keys/{{key}}")
	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetServiceAccountKeyApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: ServiceAccountKeyAssetType,
			Resource: &cai.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://iam.googleapis.com/$discovery/rest?version=v1",
				DiscoveryName:        "ServiceAccountKey",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}

func GetServiceAccountKeyApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	idProp, err := expandServiceAccountKeyId(d.Get("id"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("id"); !tpgresource.IsEmptyValue(reflect.ValueOf(idProp)) && (ok || !reflect.DeepEqual(v, idProp)) {
		obj["id"] = idProp
	}

	nameProp, err := expandServiceAccountKeyName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}

	privateKeyTypeProp, err := expandServiceAccountKeyPrivateKeyType(d.Get("privateKeyType"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("privateKeyType"); !tpgresource.IsEmptyValue(reflect.ValueOf(privateKeyTypeProp)) && (ok || !reflect.DeepEqual(v, privateKeyTypeProp)) {
		obj["privateKeyType"] = privateKeyTypeProp
	}

	keyAlgorithmProp, err := expandServiceAccountKeyKeyAlgorithm(d.Get("keyAlgorithm"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("keyAlgorithm"); !tpgresource.IsEmptyValue(reflect.ValueOf(keyAlgorithmProp)) && (ok || !reflect.DeepEqual(v, keyAlgorithmProp)) {
		obj["keyAlgorithm"] = keyAlgorithmProp
	}

	privateKeyDataProp, err := expandServiceAccountKeyKeyPrivateKeyData(d.Get("privateKeyData"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("privateKeyData"); !tpgresource.IsEmptyValue(reflect.ValueOf(privateKeyDataProp)) && (ok || !reflect.DeepEqual(v, privateKeyDataProp)) {
		obj["privateKeyData"] = privateKeyDataProp
	}

	publicKeyDataProp, err := expandServiceAccountKeyPublicKeyData(d.Get("publicKeyData"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("publicKeyData"); !tpgresource.IsEmptyValue(reflect.ValueOf(publicKeyDataProp)) && (ok || !reflect.DeepEqual(v, publicKeyDataProp)) {
		obj["publicKeyData"] = publicKeyDataProp
	}

	validAfterTimeProp, err := expandServiceAccountKeyValidAfterTime(d.Get("validAfterTime"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("validAfterTime"); !tpgresource.IsEmptyValue(reflect.ValueOf(validAfterTimeProp)) && (ok || !reflect.DeepEqual(v, validAfterTimeProp)) {
		obj["validAfterTime"] = validAfterTimeProp
	}

	validBeforeTimeProp, err := expandServiceAccountKeyValidBeforeTime(d.Get("validBeforeTime"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("validBeforeTime"); !tpgresource.IsEmptyValue(reflect.ValueOf(validBeforeTimeProp)) && (ok || !reflect.DeepEqual(v, validBeforeTimeProp)) {
		obj["validBeforeTime"] = validBeforeTimeProp
	}

	keyOriginProp, err := expandServiceAccountKeyKeyOrigin(d.Get("keyOrigin"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("keyOrigin"); !tpgresource.IsEmptyValue(reflect.ValueOf(keyOriginProp)) && (ok || !reflect.DeepEqual(v, keyOriginProp)) {
		obj["keyOrigin"] = keyOriginProp
	}

	keyTypeProp, err := expandServiceAccountKeykeyType(d.Get("keyType"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("keyType"); !tpgresource.IsEmptyValue(reflect.ValueOf(keyTypeProp)) && (ok || !reflect.DeepEqual(v, keyTypeProp)) {
		obj["keyType"] = keyTypeProp
	}

	disabledProp, err := expandServiceAccountKeyDisabled(d.Get("keyType"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("disabled"); !tpgresource.IsEmptyValue(reflect.ValueOf(disabledProp)) && (ok || !reflect.DeepEqual(v, disabledProp)) {
		obj["disabled"] = disabledProp
	}

	disableReasonProp, err := expandServiceAccountKeyDisableReason(d.Get("keyType"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("disableReason"); !tpgresource.IsEmptyValue(reflect.ValueOf(disableReasonProp)) && (ok || !reflect.DeepEqual(v, disableReasonProp)) {
		obj["disableReason"] = disableReasonProp
	}

	extendedStatusProp, err := expandServiceAccountKeyExtendedStatus(d.Get("extendedStatus"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("extendedStatus"); !tpgresource.IsEmptyValue(reflect.ValueOf(extendedStatusProp)) && (ok || !reflect.DeepEqual(v, extendedStatusProp)) {
		obj["extendedStatus"] = extendedStatusProp
	}

	contactProp, err := expandServiceAccountKeyContact(d.Get("contact"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("contact"); !tpgresource.IsEmptyValue(reflect.ValueOf(contactProp)) && (ok || !reflect.DeepEqual(v, contactProp)) {
		obj["contact"] = contactProp
	}

	descriptionProp, err := expandServiceAccountKeyDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}

	creatorProp, err := expandServiceAccountKeyCreator(d.Get("creator"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("creator"); !tpgresource.IsEmptyValue(reflect.ValueOf(creatorProp)) && (ok || !reflect.DeepEqual(v, creatorProp)) {
		obj["creator"] = creatorProp
	}

	return obj, nil
}

func expandServiceAccountKeyId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandServiceAccountKeyName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandServiceAccountKeyPrivateKeyType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
func expandServiceAccountKeyKeyAlgorithm(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
func expandServiceAccountKeyKeyPrivateKeyData(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
func expandServiceAccountKeyPublicKeyData(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
func expandServiceAccountKeyValidAfterTime(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
func expandServiceAccountKeyValidBeforeTime(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
func expandServiceAccountKeyKeyOrigin(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
func expandServiceAccountKeykeyType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
func expandServiceAccountKeyDisabled(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
func expandServiceAccountKeyDisableReason(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
func expandServiceAccountKeyExtendedStatus(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
func expandServiceAccountKeyContact(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
func expandServiceAccountKeyDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
func expandServiceAccountKeyCreator(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
