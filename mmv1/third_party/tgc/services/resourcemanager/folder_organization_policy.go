package resourcemanager

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

func ResourceConverterFolderOrgPolicy() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Folder",
		Convert:           GetFolderOrgPolicyCaiObject,
		MergeCreateUpdate: MergeFolderOrgPolicy,
	}
}

func GetFolderOrgPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//cloudresourcemanager.googleapis.com/{{folder}}")
	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetFolderOrgPolicyApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name:      name,
			Type:      "cloudresourcemanager.googleapis.com/Folder",
			OrgPolicy: []*cai.OrgPolicy{&obj},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}

func MergeFolderOrgPolicy(existing, incoming cai.Asset) cai.Asset {
	existing.OrgPolicy = append(existing.OrgPolicy, incoming.OrgPolicy...)
	return existing
}

func GetFolderOrgPolicyApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (cai.OrgPolicy, error) {

	listPolicy, err := expandListOrganizationPolicy(d.Get("list_policy").([]interface{}))
	if err != nil {
		return cai.OrgPolicy{}, err
	}

	restoreDefault, err := expandRestoreOrganizationPolicy(d.Get("restore_policy").([]interface{}))
	if err != nil {
		return cai.OrgPolicy{}, err
	}

	policy := cai.OrgPolicy{
		Constraint:     CanonicalOrgPolicyConstraint(d.Get("constraint").(string)),
		BooleanPolicy:  expandBooleanOrganizationPolicy(d.Get("boolean_policy").([]interface{})),
		ListPolicy:     listPolicy,
		RestoreDefault: restoreDefault,
	}

	return policy, nil
}
