package firebase

import (
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const firebaseProjectAssetType string = "firebase.googleapis.com/FirebaseProject"

func ResourceConverterFirebaseProject() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: firebaseProjectAssetType,
		Convert:   GetFirebaseProjectCaiObject,
	}
}

func GetFirebaseProjectCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//firebase.googleapis.com/v1beta1/projects/{{project}}")
	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetFirebaseProjectApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: firebaseProjectAssetType,
			Resource: &cai.AssetResource{
				Version:              "v1beta1",
				DiscoveryDocumentURI: "https://firebase.googleapis.com/$discovery/rest?version=v1beta1",
				DiscoveryName:        "FirebaseProject",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}

func GetFirebaseProjectApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	nameProp, err := expandFirebaseProjectName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}

	projectProp, err := expandFirebaseProjectProjectId(d.Get("project"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("project"); !tpgresource.IsEmptyValue(reflect.ValueOf(projectProp)) && (ok || !reflect.DeepEqual(v, projectProp)) {
		obj["projectId"] = projectProp
	}

	idProp, err := expandFirebaseProjectId(d.Get("id"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("id"); !tpgresource.IsEmptyValue(reflect.ValueOf(idProp)) && (ok || !reflect.DeepEqual(v, idProp)) {
		obj["id"] = idProp
	}

	projectNumberProp, err := expandFirebaseProjectProjectNumber(d.Get("project_number"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("project_number"); !tpgresource.IsEmptyValue(reflect.ValueOf(projectNumberProp)) && (ok || !reflect.DeepEqual(v, projectNumberProp)) {
		obj["projectNumber"] = projectNumberProp
	}

	displayNameProp, err := expandFirebaseProjectDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("display_name"); !tpgresource.IsEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}

	return obj, nil
}

func expandFirebaseProjectDisplayName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandFirebaseProjectProjectNumber(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandFirebaseProjectId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	v, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}")
	if err != nil {
		return nil, err
	}

	return v, nil
}

func expandFirebaseProjectProjectId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandFirebaseProjectName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
