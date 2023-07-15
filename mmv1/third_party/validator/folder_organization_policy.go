package google

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
)

func resourceConverterFolderOrgPolicy() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Folder",
		Convert:           GetFolderOrgPolicyCaiObject,
		MergeCreateUpdate: MergeFolderOrgPolicy,
	}
}

func GetFolderOrgPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	name, err := tpgresource.AssetName(d, config, "//cloudresourcemanager.googleapis.com/{{folder}}")
	if err != nil {
		return []tpgresource.Asset{}, err
	}
	if obj, err := GetFolderOrgPolicyApiObject(d, config); err == nil {
		return []tpgresource.Asset{{
			Name:      name,
			Type:      "cloudresourcemanager.googleapis.com/Folder",
			OrgPolicy: []*tpgresource.OrgPolicy{&obj},
		}}, nil
	} else {
		return []tpgresource.Asset{}, err
	}
}

func MergeFolderOrgPolicy(existing, incoming tpgresource.Asset) tpgresource.Asset {
	existing.OrgPolicy = append(existing.OrgPolicy, incoming.OrgPolicy...)
	return existing
}

func GetFolderOrgPolicyApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgresource.OrgPolicy, error) {

	listPolicy, err := expandListOrganizationPolicy(d.Get("list_policy").([]interface{}))
	if err != nil {
		return tpgresource.OrgPolicy{}, err
	}

	restoreDefault, err := expandRestoreOrganizationPolicy(d.Get("restore_policy").([]interface{}))
	if err != nil {
		return tpgresource.OrgPolicy{}, err
	}

	policy := tpgresource.OrgPolicy{
		Constraint:     CanonicalOrgPolicyConstraint(d.Get("constraint").(string)),
		BooleanPolicy:  expandBooleanOrganizationPolicy(d.Get("boolean_policy").([]interface{})),
		ListPolicy:     listPolicy,
		RestoreDefault: restoreDefault,
	}

	return policy, nil
}
