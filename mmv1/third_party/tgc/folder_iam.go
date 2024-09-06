package resourcemanager

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

func ResourceConverterFolderIamPolicy() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Folder",
		Convert:           GetFolderIamPolicyCaiObject,
		MergeCreateUpdate: MergeFolderIamPolicy,
	}
}

func ResourceConverterFolderIamBinding() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Folder",
		Convert:           GetFolderIamBindingCaiObject,
		FetchFullResource: FetchFolderIamPolicy,
		MergeCreateUpdate: MergeFolderIamBinding,
		MergeDelete:       MergeFolderIamBindingDelete,
	}
}

func ResourceConverterFolderIamMember() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Folder",
		Convert:           GetFolderIamMemberCaiObject,
		FetchFullResource: FetchFolderIamPolicy,
		MergeCreateUpdate: MergeFolderIamMember,
		MergeDelete:       MergeFolderIamMemberDelete,
	}
}

func GetFolderIamPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newFolderIamAsset(d, config, cai.ExpandIamPolicyBindings)
}

func GetFolderIamBindingCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newFolderIamAsset(d, config, cai.ExpandIamRoleBindings)
}

func GetFolderIamMemberCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newFolderIamAsset(d, config, cai.ExpandIamMemberBindings)
}

func MergeFolderIamPolicy(existing, incoming cai.Asset) cai.Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeFolderIamBinding(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeIamAssets(existing, incoming, cai.MergeAuthoritativeBindings)
}

func MergeFolderIamBindingDelete(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeDeleteIamAssets(existing, incoming, cai.MergeDeleteAuthoritativeBindings)
}

func MergeFolderIamMember(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeIamAssets(existing, incoming, cai.MergeAdditiveBindings)
}

func MergeFolderIamMemberDelete(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeDeleteIamAssets(existing, incoming, cai.MergeDeleteAdditiveBindings)
}

func newFolderIamAsset(
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	expandBindings func(d tpgresource.TerraformResourceData) ([]cai.IAMBinding, error),
) ([]cai.Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []cai.Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	// The "folder" argument is of the form "folders/12345"
	name, err := cai.AssetName(d, config, "//cloudresourcemanager.googleapis.com/{{folder}}")
	if err != nil {
		return []cai.Asset{}, err
	}

	return []cai.Asset{{
		Name: name,
		Type: "cloudresourcemanager.googleapis.com/Folder",
		IAMPolicy: &cai.IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchFolderIamPolicy(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (cai.Asset, error) {
	if _, ok := d.GetOk("folder"); !ok {
		return cai.Asset{}, cai.ErrEmptyIdentityField
	}

	return cai.FetchIamPolicy(
		NewFolderIamUpdater,
		d,
		config,
		"//cloudresourcemanager.googleapis.com/{{folder}}",
		"cloudresourcemanager.googleapis.com/Folder",
	)
}
