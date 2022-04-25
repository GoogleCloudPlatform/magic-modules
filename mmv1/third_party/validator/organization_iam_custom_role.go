package google

import (
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const OrganizationIAMCustomRoleAssetType string = "iam.googleapis.com/Role"

func resourceConverterOrganizationIAMCustomRole() ResourceConverter {
	return ResourceConverter{
		AssetType: OrganizationIAMCustomRoleAssetType,
		Convert:   GetOrganizationIAMCustomRoleCaiObject,
	}
}

func GetOrganizationIAMCustomRoleCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//iam.googleapis.com/organizations/{{org_id}}/roles/{{role_id}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetOrganizationIAMCustomRoleApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: OrganizationIAMCustomRoleAssetType,
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

func GetOrganizationIAMCustomRoleApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	descriptionProp, err := expandOrganizationIAMCustomRoleDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}

	titleProp, err := expandOrganizationIAMCustomRoleTitle(d.Get("title"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("title"); !isEmptyValue(reflect.ValueOf(titleProp)) && (ok || !reflect.DeepEqual(v, titleProp)) {
		obj["title"] = titleProp
	}

	stageProp, err := expandOrganizationIAMCustomRoleStage(d.Get("stage"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("stage"); !isEmptyValue(reflect.ValueOf(stageProp)) && (ok || !reflect.DeepEqual(v, stageProp)) {
		obj["stage"] = stageProp
	}

	includedPermissionsProp, err := expandOrganizationIAMCustomRolePermissions(d.Get("permissions"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("permissions"); !isEmptyValue(reflect.ValueOf(includedPermissionsProp)) && (ok || !reflect.DeepEqual(v, includedPermissionsProp)) {
		obj["includedPermissions"] = includedPermissionsProp
	}

	return obj, nil
}

func expandOrganizationIAMCustomRoleDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandOrganizationIAMCustomRoleTitle(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandOrganizationIAMCustomRoleStage(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandOrganizationIAMCustomRolePermissions(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v.(*schema.Set).List(), nil
}
