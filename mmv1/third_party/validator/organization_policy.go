package google

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
)

func resourceConverterOrganizationPolicy() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Organization",
		Convert:           GetOrganizationPolicyCaiObject,
		MergeCreateUpdate: MergeOrganizationPolicy,
	}
}

func GetOrganizationPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	name, err := tpgresource.AssetName(d, config, "//cloudresourcemanager.googleapis.com/organizations/{{org_id}}")
	if err != nil {
		return []tpgresource.Asset{}, err
	}
	if obj, err := GetOrganizationPolicyApiObject(d, config); err == nil {
		return []tpgresource.Asset{{
			Name:      name,
			Type:      "cloudresourcemanager.googleapis.com/Organization",
			OrgPolicy: []*tpgresource.OrgPolicy{&obj},
		}}, nil
	} else {
		return []tpgresource.Asset{}, err
	}
}

func MergeOrganizationPolicy(existing, incoming tpgresource.Asset) tpgresource.Asset {
	existing.OrgPolicy = append(existing.OrgPolicy, incoming.OrgPolicy...)
	return existing
}

func GetOrganizationPolicyApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgresource.OrgPolicy, error) {

	listPolicy, err := expandListOrganizationPolicy(d.Get("list_policy").([]interface{}))
	if err != nil {
		return tpgresource.OrgPolicy{}, err
	}

	restoreDefault, err := expandRestoreOrganizationPolicy(d.Get("restore_policy").([]interface{}))
	if err != nil {
		return tpgresource.OrgPolicy{}, err
	}

	policy := tpgresource.OrgPolicy{
		Constraint:     canonicalOrgPolicyConstraint(d.Get("constraint").(string)),
		BooleanPolicy:  expandBooleanOrganizationPolicy(d.Get("boolean_policy").([]interface{})),
		ListPolicy:     listPolicy,
		RestoreDefault: restoreDefault,
	}

	return policy, nil
}
