package Apikeys

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)


const ApikeysKeyAssetType string = "apikeys.googleapis.com/Key"

func GetApikeysKeyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//apikeys.googleapis.com/v2/projects/{{project}}/locations/global/keys/{{key}})
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
				DiscoveryName:        "key",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}


func GetApikeysKeyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	obj := make(map[string]interface{})
	
 	nameProp, err := expandApikeysKeyName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	
	uidProp, err := expandApikeysKeyUid(d.Get("uid"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("uid"); !tpgresource.IsEmptyValue(reflect.ValueOf(uidProp)) && (ok || !reflect.DeepEqual(v, uidProp)) {
		obj["uid"] = nameProp
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
	
	restrictionsProp, err := expandApikeysKeyDRestrictions(d.Get("restrictions"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("restrictions"); !tpgresource.IsEmptyValue(reflect.ValueOf(restrictionsProp)) && (ok || !reflect.DeepEqual(v, restrictionsProp)) {
		obj["restrictions"] = restrictionsProp
	}
	
	etagProp, err := expandApikeysKeyDRestrictions(d.Get("etag"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("etag"); !tpgresource.IsEmptyValue(reflect.ValueOf(etagProp)) && (ok || !reflect.DeepEqual(v, etagProp)) {
		obj["etag"] = etagProp
	}
	
	
	
	
	return obj, nil
}





