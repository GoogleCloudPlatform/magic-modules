package resourcemanager

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

func ResourceConverterProjectIamPolicy() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Project",
		Convert:           GetProjectIamPolicyCaiObject,
		MergeCreateUpdate: MergeProjectIamPolicy,
	}
}

func ResourceConverterProjectIamBinding() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Project",
		Convert:           GetProjectIamBindingCaiObject,
		FetchFullResource: FetchProjectIamPolicy,
		MergeCreateUpdate: MergeProjectIamBinding,
		MergeDelete:       MergeProjectIamBindingDelete,
	}
}

func ResourceConverterProjectIamMember() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Project",
		Convert:           GetProjectIamMemberCaiObject,
		FetchFullResource: FetchProjectIamPolicy,
		MergeCreateUpdate: MergeProjectIamMember,
		MergeDelete:       MergeProjectIamMemberDelete,
	}
}

func GetProjectIamPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newProjectIamAsset(d, config, cai.ExpandIamPolicyBindings)
}

func GetProjectIamBindingCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newProjectIamAsset(d, config, cai.ExpandIamRoleBindings)
}

func GetProjectIamMemberCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newProjectIamAsset(d, config, cai.ExpandIamMemberBindings)
}

func MergeProjectIamPolicy(existing, incoming cai.Asset) cai.Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeProjectIamBinding(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeIamAssets(existing, incoming, cai.MergeAuthoritativeBindings)
}

func MergeProjectIamBindingDelete(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeDeleteIamAssets(existing, incoming, cai.MergeDeleteAuthoritativeBindings)
}

func MergeProjectIamMember(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeIamAssets(existing, incoming, cai.MergeAdditiveBindings)
}

func MergeProjectIamMemberDelete(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeDeleteIamAssets(existing, incoming, cai.MergeDeleteAdditiveBindings)
}

func newProjectIamAsset(
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	expandBindings func(d tpgresource.TerraformResourceData) ([]cai.IAMBinding, error),
) ([]cai.Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []cai.Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	// Ideally we should use project_number, but since that is generated server-side,
	// we substitute project_id.
	name, err := cai.AssetName(d, config, "//cloudresourcemanager.googleapis.com/projects/{{project}}")
	if err != nil {
		return []cai.Asset{}, err
	}

	return []cai.Asset{{
		Name: name,
		Type: "cloudresourcemanager.googleapis.com/Project",
		IAMPolicy: &cai.IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchProjectIamPolicy(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (cai.Asset, error) {
	if _, ok := d.GetOk("project"); !ok {
		return cai.Asset{}, cai.ErrEmptyIdentityField
	}

	// We use project_id in the asset name template to be consistent with newProjectIamAsset.
	return cai.FetchIamPolicy(
		NewProjectIamUpdater,
		d,
		config,
		"//cloudresourcemanager.googleapis.com/projects/{{project}}",
		"cloudresourcemanager.googleapis.com/Project",
	)
}
