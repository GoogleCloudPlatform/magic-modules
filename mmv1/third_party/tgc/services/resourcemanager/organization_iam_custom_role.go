package resourcemanager

import (
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const OrganizationIAMCustomRoleAssetType string = "iam.googleapis.com/Role"

func ResourceConverterOrganizationIAMCustomRole() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: OrganizationIAMCustomRoleAssetType,
		Convert:   GetOrganizationIAMCustomRoleCaiObject,
	}
}

func GetOrganizationIAMCustomRoleCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//iam.googleapis.com/organizations/{{org_id}}/roles/{{role_id}}")
	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetOrganizationIAMCustomRoleApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: OrganizationIAMCustomRoleAssetType,
			Resource: &cai.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://iam.googleapis.com/$discovery/rest?version=v1",
				DiscoveryName:        "Role",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}

func GetOrganizationIAMCustomRoleApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	descriptionProp, err := expandOrganizationIAMCustomRoleDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}

	titleProp, err := expandOrganizationIAMCustomRoleTitle(d.Get("title"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("title"); !tpgresource.IsEmptyValue(reflect.ValueOf(titleProp)) && (ok || !reflect.DeepEqual(v, titleProp)) {
		obj["title"] = titleProp
	}

	stageProp, err := expandOrganizationIAMCustomRoleStage(d.Get("stage"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("stage"); !tpgresource.IsEmptyValue(reflect.ValueOf(stageProp)) && (ok || !reflect.DeepEqual(v, stageProp)) {
		obj["stage"] = stageProp
	}

	includedPermissionsProp, err := expandOrganizationIAMCustomRolePermissions(d.Get("permissions"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("permissions"); !tpgresource.IsEmptyValue(reflect.ValueOf(includedPermissionsProp)) && (ok || !reflect.DeepEqual(v, includedPermissionsProp)) {
		obj["includedPermissions"] = includedPermissionsProp
	}

	return obj, nil
}

func expandOrganizationIAMCustomRoleDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOrganizationIAMCustomRoleTitle(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOrganizationIAMCustomRoleStage(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOrganizationIAMCustomRolePermissions(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v.(*schema.Set).List(), nil
}
