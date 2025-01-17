package appengine

import (
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const AppEngineApplicationAssetType string = "appengine.googleapis.com/Application"

func ResourceConverterAppEngineApplication() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: AppEngineApplicationAssetType,
		Convert:   GetAppEngineApplicationCaiObject,
	}
}

func GetAppEngineApplicationCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//appengine.googleapis.com/v1/{{name}}")
	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetAppEngineApplicationApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: AppEngineApplicationAssetType,
			Resource: &cai.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/appengine/v1beta/rest",
				DiscoveryName:        "Application",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}

func GetAppEngineApplicationApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	idProp, err := expandAppEngineApplicationId(d.Get("id"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("id"); !tpgresource.IsEmptyValue(reflect.ValueOf(idProp)) && (ok || !reflect.DeepEqual(v, idProp)) {
		obj["id"] = idProp
	}

	locationIdProp, err := expandAppEngineApplicationLocationId(d.Get("location_id"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("locationId"); !tpgresource.IsEmptyValue(reflect.ValueOf(locationIdProp)) && (ok || !reflect.DeepEqual(v, locationIdProp)) {
		obj["location_id"] = locationIdProp
	}

	return obj, nil
}

func expandAppEngineApplicationId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineApplicationLocationId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
