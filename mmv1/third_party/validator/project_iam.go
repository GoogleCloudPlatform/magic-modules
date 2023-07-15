package google

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgiamresource"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
)

func resourceConverterProjectIamPolicy() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Project",
		Convert:           GetProjectIamPolicyCaiObject,
		MergeCreateUpdate: MergeProjectIamPolicy,
	}
}

func resourceConverterProjectIamBinding() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Project",
		Convert:           GetProjectIamBindingCaiObject,
		FetchFullResource: FetchProjectIamPolicy,
		MergeCreateUpdate: MergeProjectIamBinding,
		MergeDelete:       MergeProjectIamBindingDelete,
	}
}

func resourceConverterProjectIamMember() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Project",
		Convert:           GetProjectIamMemberCaiObject,
		FetchFullResource: FetchProjectIamPolicy,
		MergeCreateUpdate: MergeProjectIamMember,
		MergeDelete:       MergeProjectIamMemberDelete,
	}
}

func GetProjectIamPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newProjectIamAsset(d, config, tpgiamresource.ExpandIamPolicyBindings)
}

func GetProjectIamBindingCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newProjectIamAsset(d, config, tpgiamresource.ExpandIamRoleBindings)
}

func GetProjectIamMemberCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newProjectIamAsset(d, config, tpgiamresource.ExpandIamMemberBindings)
}

func MergeProjectIamPolicy(existing, incoming tpgresource.Asset) tpgresource.Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeProjectIamBinding(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeIamAssets(existing, incoming, tpgiamresource.MergeAuthoritativeBindings)
}

func MergeProjectIamBindingDelete(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeDeleteIamAssets(existing, incoming, tpgiamresource.MergeDeleteAuthoritativeBindings)
}

func MergeProjectIamMember(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeIamAssets(existing, incoming, tpgiamresource.MergeAdditiveBindings)
}

func MergeProjectIamMemberDelete(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeDeleteIamAssets(existing, incoming, tpgiamresource.MergeDeleteAdditiveBindings)
}

func newProjectIamAsset(
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	expandBindings func(d tpgresource.TerraformResourceData) ([]tpgresource.IAMBinding, error),
) ([]tpgresource.Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []tpgresource.Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	// Ideally we should use project_number, but since that is generated server-side,
	// we substitute project_id.
	name, err := tpgresource.AssetName(d, config, "//cloudresourcemanager.googleapis.com/projects/{{project}}")
	if err != nil {
		return []tpgresource.Asset{}, err
	}

	return []tpgresource.Asset{{
		Name: name,
		Type: "cloudresourcemanager.googleapis.com/Project",
		IAMPolicy: &tpgresource.IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchProjectIamPolicy(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgresource.Asset, error) {
	if _, ok := d.GetOk("project"); !ok {
		return tpgresource.Asset{}, tpgresource.ErrEmptyIdentityField
	}

	// We use project_id in the asset name template to be consistent with newProjectIamAsset.
	return tpgiamresource.FetchIamPolicy(
		NewProjectIamUpdater,
		d,
		config,
		"//cloudresourcemanager.googleapis.com/projects/{{project}}",
		"cloudresourcemanager.googleapis.com/Project",
	)
}
