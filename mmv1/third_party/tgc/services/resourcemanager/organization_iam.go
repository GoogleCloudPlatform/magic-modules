package resourcemanager

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

func ResourceConverterOrganizationIamPolicy() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Organization",
		Convert:           GetOrganizationIamPolicyCaiObject,
		MergeCreateUpdate: MergeOrganizationIamPolicy,
	}
}

func ResourceConverterOrganizationIamBinding() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Organization",
		Convert:           GetOrganizationIamBindingCaiObject,
		FetchFullResource: FetchOrganizationIamPolicy,
		MergeCreateUpdate: MergeOrganizationIamBinding,
		MergeDelete:       MergeOrganizationIamBindingDelete,
	}
}

func ResourceConverterOrganizationIamMember() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Organization",
		Convert:           GetOrganizationIamMemberCaiObject,
		FetchFullResource: FetchOrganizationIamPolicy,
		MergeCreateUpdate: MergeOrganizationIamMember,
		MergeDelete:       MergeOrganizationIamMemberDelete,
	}
}

func GetOrganizationIamPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newOrganizationIamAsset(d, config, cai.ExpandIamPolicyBindings)
}

func GetOrganizationIamBindingCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newOrganizationIamAsset(d, config, cai.ExpandIamRoleBindings)
}

func GetOrganizationIamMemberCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newOrganizationIamAsset(d, config, cai.ExpandIamMemberBindings)
}

func MergeOrganizationIamPolicy(existing, incoming cai.Asset) cai.Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeOrganizationIamBinding(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeIamAssets(existing, incoming, cai.MergeAuthoritativeBindings)
}

func MergeOrganizationIamBindingDelete(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeDeleteIamAssets(existing, incoming, cai.MergeDeleteAuthoritativeBindings)
}

func MergeOrganizationIamMember(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeIamAssets(existing, incoming, cai.MergeAdditiveBindings)
}

func MergeOrganizationIamMemberDelete(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeDeleteIamAssets(existing, incoming, cai.MergeDeleteAdditiveBindings)
}

func newOrganizationIamAsset(
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	expandBindings func(d tpgresource.TerraformResourceData) ([]cai.IAMBinding, error),
) ([]cai.Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []cai.Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := cai.AssetName(d, config, "//cloudresourcemanager.googleapis.com/organizations/{{org_id}}")
	if err != nil {
		return []cai.Asset{}, err
	}

	return []cai.Asset{{
		Name: name,
		Type: "cloudresourcemanager.googleapis.com/Organization",
		IAMPolicy: &cai.IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchOrganizationIamPolicy(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (cai.Asset, error) {
	return cai.FetchIamPolicy(
		NewOrganizationIamUpdater,
		d,
		config,
		"//cloudresourcemanager.googleapis.com/organizations/{{org_id}}",
		"cloudresourcemanager.googleapis.com/Organization",
	)
}
