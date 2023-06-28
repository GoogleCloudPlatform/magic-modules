package google

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgiamresource"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
)

// Provide a separate asset type constant so we don't have to worry about name conflicts between IAM and non-IAM converter files
const StorageBucketIAMAssetType string = "storage.googleapis.com/Bucket"

func resourceConverterStorageBucketIamPolicy() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         StorageBucketIAMAssetType,
		Convert:           GetStorageBucketIamPolicyCaiObject,
		MergeCreateUpdate: MergeStorageBucketIamPolicy,
	}
}

func resourceConverterStorageBucketIamBinding() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         StorageBucketIAMAssetType,
		Convert:           GetStorageBucketIamBindingCaiObject,
		FetchFullResource: FetchStorageBucketIamPolicy,
		MergeCreateUpdate: MergeStorageBucketIamBinding,
		MergeDelete:       MergeStorageBucketIamBindingDelete,
	}
}

func resourceConverterStorageBucketIamMember() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         StorageBucketIAMAssetType,
		Convert:           GetStorageBucketIamMemberCaiObject,
		FetchFullResource: FetchStorageBucketIamPolicy,
		MergeCreateUpdate: MergeStorageBucketIamMember,
		MergeDelete:       MergeStorageBucketIamMemberDelete,
	}
}

func GetStorageBucketIamPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newStorageBucketIamAsset(d, config, tpgiamresource.ExpandIamPolicyBindings)
}

func GetStorageBucketIamBindingCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newStorageBucketIamAsset(d, config, tpgiamresource.ExpandIamRoleBindings)
}

func GetStorageBucketIamMemberCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newStorageBucketIamAsset(d, config, tpgiamresource.ExpandIamMemberBindings)
}

func MergeStorageBucketIamPolicy(existing, incoming tpgresource.Asset) tpgresource.Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeStorageBucketIamBinding(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeIamAssets(existing, incoming, tpgiamresource.MergeAuthoritativeBindings)
}

func MergeStorageBucketIamBindingDelete(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeDeleteIamAssets(existing, incoming, tpgiamresource.MergeDeleteAuthoritativeBindings)
}

func MergeStorageBucketIamMember(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeIamAssets(existing, incoming, tpgiamresource.MergeAdditiveBindings)
}

func MergeStorageBucketIamMemberDelete(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeDeleteIamAssets(existing, incoming, tpgiamresource.MergeDeleteAdditiveBindings)
}

func newStorageBucketIamAsset(
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	expandBindings func(d tpgresource.TerraformResourceData) ([]tpgresource.IAMBinding, error),
) ([]tpgresource.Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []tpgresource.Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := tpgresource.AssetName(d, config, "//storage.googleapis.com/{{bucket}}")
	if err != nil {
		return []tpgresource.Asset{}, err
	}

	return []tpgresource.Asset{{
		Name: name,
		Type: StorageBucketIAMAssetType,
		IAMPolicy: &tpgresource.IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchStorageBucketIamPolicy(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgresource.Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("bucket"); !ok {
		return tpgresource.Asset{}, tpgresource.ErrEmptyIdentityField
	}

	return tpgiamresource.FetchIamPolicy(
		StorageBucketIamUpdaterProducer,
		d,
		config,
		"//storage.googleapis.com/{{bucket}}",
		StorageBucketIAMAssetType,
	)
}
