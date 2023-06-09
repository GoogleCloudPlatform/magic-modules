package google

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgiamresource"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
)

func resourceConverterFolderIamPolicy() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Folder",
		Convert:           GetFolderIamPolicyCaiObject,
		MergeCreateUpdate: MergeFolderIamPolicy,
	}
}

func resourceConverterFolderIamBinding() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Folder",
		Convert:           GetFolderIamBindingCaiObject,
		FetchFullResource: FetchFolderIamPolicy,
		MergeCreateUpdate: MergeFolderIamBinding,
		MergeDelete:       MergeFolderIamBindingDelete,
	}
}

func resourceConverterFolderIamMember() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Folder",
		Convert:           GetFolderIamMemberCaiObject,
		FetchFullResource: FetchFolderIamPolicy,
		MergeCreateUpdate: MergeFolderIamMember,
		MergeDelete:       MergeFolderIamMemberDelete,
	}
}

func GetFolderIamPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newFolderIamAsset(d, config, tpgiamresource.ExpandIamPolicyBindings)
}

func GetFolderIamBindingCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newFolderIamAsset(d, config, tpgiamresource.ExpandIamRoleBindings)
}

func GetFolderIamMemberCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newFolderIamAsset(d, config, tpgiamresource.ExpandIamMemberBindings)
}

func MergeFolderIamPolicy(existing, incoming tpgresource.Asset) tpgresource.Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeFolderIamBinding(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeIamAssets(existing, incoming, tpgiamresource.MergeAuthoritativeBindings)
}

func MergeFolderIamBindingDelete(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeDeleteIamAssets(existing, incoming, tpgiamresource.MergeDeleteAuthoritativeBindings)
}

func MergeFolderIamMember(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeIamAssets(existing, incoming, tpgiamresource.MergeAdditiveBindings)
}

func MergeFolderIamMemberDelete(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeDeleteIamAssets(existing, incoming, tpgiamresource.MergeDeleteAdditiveBindings)
}

func newFolderIamAsset(
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	expandBindings func(d tpgresource.TerraformResourceData) ([]tpgresource.IAMBinding, error),
) ([]tpgresource.Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []tpgresource.Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	// The "folder" argument is of the form "folders/12345"
	name, err := tpgresource.AssetName(d, config, "//cloudresourcemanager.googleapis.com/{{folder}}")
	if err != nil {
		return []tpgresource.Asset{}, err
	}

	return []tpgresource.Asset{{
		Name: name,
		Type: "cloudresourcemanager.googleapis.com/Folder",
		IAMPolicy: &tpgresource.IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchFolderIamPolicy(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgresource.Asset, error) {
	if _, ok := d.GetOk("folder"); !ok {
		return tpgresource.Asset{}, tpgresource.ErrEmptyIdentityField
	}

	return tpgiamresource.FetchIamPolicy(
		NewFolderIamUpdater,
		d,
		config,
		"//cloudresourcemanager.googleapis.com/{{folder}}",
		"cloudresourcemanager.googleapis.com/Folder",
	)
}
