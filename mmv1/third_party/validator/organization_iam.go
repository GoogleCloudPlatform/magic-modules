package google

import "fmt"

func resourceConverterOrganizationIamPolicy() ResourceConverter {
	return ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Organization",
		Convert:           GetOrganizationIamPolicyCaiObject,
		MergeCreateUpdate: MergeOrganizationIamPolicy,
	}
}

func resourceConverterOrganizationIamBinding() ResourceConverter {
	return ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Organization",
		Convert:           GetOrganizationIamBindingCaiObject,
		FetchFullResource: FetchOrganizationIamPolicy,
		MergeCreateUpdate: MergeOrganizationIamBinding,
		MergeDelete:       MergeOrganizationIamBindingDelete,
	}
}

func resourceConverterOrganizationIamMember() ResourceConverter {
	return ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Organization",
		Convert:           GetOrganizationIamMemberCaiObject,
		FetchFullResource: FetchOrganizationIamPolicy,
		MergeCreateUpdate: MergeOrganizationIamMember,
		MergeDelete:       MergeOrganizationIamMemberDelete,
	}
}

func GetOrganizationIamPolicyCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newOrganizationIamAsset(d, config, expandIamPolicyBindings)
}

func GetOrganizationIamBindingCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newOrganizationIamAsset(d, config, expandIamRoleBindings)
}

func GetOrganizationIamMemberCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newOrganizationIamAsset(d, config, expandIamMemberBindings)
}

func MergeOrganizationIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeOrganizationIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeOrganizationIamBindingDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAuthoritativeBindings)
}

func MergeOrganizationIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func MergeOrganizationIamMemberDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAdditiveBindings)
}

func newOrganizationIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) ([]Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := assetName(d, config, "//cloudresourcemanager.googleapis.com/organizations/{{org_id}}")
	if err != nil {
		return []Asset{}, err
	}

	return []Asset{{
		Name: name,
		Type: "cloudresourcemanager.googleapis.com/Organization",
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchOrganizationIamPolicy(d TerraformResourceData, config *Config) (Asset, error) {
	return fetchIamPolicy(
		NewOrganizationIamUpdater,
		d,
		config,
		"//cloudresourcemanager.googleapis.com/organizations/{{org_id}}",
		"cloudresourcemanager.googleapis.com/Organization",
	)
}
