package google

import (
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ProjectIAMCustomRoleAssetType string = "iam.googleapis.com/Role"

func resourceConverterProjectIAMCustomRole() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType: ProjectIAMCustomRoleAssetType,
		Convert:   GetProjectIAMCustomRoleCaiObject,
	}
}

func GetProjectIAMCustomRoleCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	name, err := tpgresource.AssetName(d, config, "//iam.googleapis.com/projects/{{project}}/roles/{{role_id}}")
	if err != nil {
		return []tpgresource.Asset{}, err
	}
	if obj, err := GetProjectIAMCustomRoleApiObject(d, config); err == nil {
		return []tpgresource.Asset{{
			Name: name,
			Type: ProjectIAMCustomRoleAssetType,
			Resource: &tpgresource.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://iam.googleapis.com/$discovery/rest?version=v1",
				DiscoveryName:        "Role",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []tpgresource.Asset{}, err
	}
}

func GetProjectIAMCustomRoleApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	descriptionProp, err := expandProjectIAMCustomRoleDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}

	titleProp, err := expandProjectIAMCustomRoleTitle(d.Get("title"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("title"); !tpgresource.IsEmptyValue(reflect.ValueOf(titleProp)) && (ok || !reflect.DeepEqual(v, titleProp)) {
		obj["title"] = titleProp
	}

	stageProp, err := expandProjectIAMCustomRoleStage(d.Get("stage"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("stage"); !tpgresource.IsEmptyValue(reflect.ValueOf(stageProp)) && (ok || !reflect.DeepEqual(v, stageProp)) {
		obj["stage"] = stageProp
	}

	includedPermissionsProp, err := expandProjectIAMCustomRolePermissions(d.Get("permissions"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("permissions"); !tpgresource.IsEmptyValue(reflect.ValueOf(includedPermissionsProp)) && (ok || !reflect.DeepEqual(v, includedPermissionsProp)) {
		obj["includedPermissions"] = includedPermissionsProp
	}

	return obj, nil
}

func expandProjectIAMCustomRoleDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandProjectIAMCustomRoleTitle(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandProjectIAMCustomRoleStage(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandProjectIAMCustomRolePermissions(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v.(*schema.Set).List(), nil
}
