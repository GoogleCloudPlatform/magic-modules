package google

import (
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ProjectIAMCustomRoleAssetType string = "iam.googleapis.com/Role"

func resourceConverterProjectIAMCustomRole() ResourceConverter {
	return ResourceConverter{
		AssetType: ProjectIAMCustomRoleAssetType,
		Convert:   GetProjectIAMCustomRoleCaiObject,
	}
}

func GetProjectIAMCustomRoleCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//iam.googleapis.com/projects/{{project}}/roles/{{role_id}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetProjectIAMCustomRoleApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: ProjectIAMCustomRoleAssetType,
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://iam.googleapis.com/$discovery/rest?version=v1",
				DiscoveryName:        "Role",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetProjectIAMCustomRoleApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	descriptionProp, err := expandProjectIAMCustomRoleDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}

	titleProp, err := expandProjectIAMCustomRoleTitle(d.Get("title"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("title"); !isEmptyValue(reflect.ValueOf(titleProp)) && (ok || !reflect.DeepEqual(v, titleProp)) {
		obj["title"] = titleProp
	}

	stageProp, err := expandProjectIAMCustomRoleStage(d.Get("stage"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("stage"); !isEmptyValue(reflect.ValueOf(stageProp)) && (ok || !reflect.DeepEqual(v, stageProp)) {
		obj["stage"] = stageProp
	}

	includedPermissionsProp, err := expandProjectIAMCustomRolePermissions(d.Get("permissions"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("permissions"); !isEmptyValue(reflect.ValueOf(includedPermissionsProp)) && (ok || !reflect.DeepEqual(v, includedPermissionsProp)) {
		obj["includedPermissions"] = includedPermissionsProp
	}

	return obj, nil
}

func expandProjectIAMCustomRoleDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandProjectIAMCustomRoleTitle(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandProjectIAMCustomRoleStage(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandProjectIAMCustomRolePermissions(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v.(*schema.Set).List(), nil
}
