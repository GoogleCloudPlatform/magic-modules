package google

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
)

// Provide a separate asset type constant so we don't have to worry about name conflicts between IAM and non-IAM converter files
const StorageBucketIAMAssetType string = "storage.googleapis.com/Bucket"

func resourceConverterStorageBucketIamPolicy() ResourceConverter {
	return ResourceConverter{
		AssetType:         StorageBucketIAMAssetType,
		Convert:           GetStorageBucketIamPolicyCaiObject,
		MergeCreateUpdate: MergeStorageBucketIamPolicy,
	}
}

func resourceConverterStorageBucketIamBinding() ResourceConverter {
	return ResourceConverter{
		AssetType:         StorageBucketIAMAssetType,
		Convert:           GetStorageBucketIamBindingCaiObject,
		FetchFullResource: FetchStorageBucketIamPolicy,
		MergeCreateUpdate: MergeStorageBucketIamBinding,
		MergeDelete:       MergeStorageBucketIamBindingDelete,
	}
}

func resourceConverterStorageBucketIamMember() ResourceConverter {
	return ResourceConverter{
		AssetType:         StorageBucketIAMAssetType,
		Convert:           GetStorageBucketIamMemberCaiObject,
		FetchFullResource: FetchStorageBucketIamPolicy,
		MergeCreateUpdate: MergeStorageBucketIamMember,
		MergeDelete:       MergeStorageBucketIamMemberDelete,
	}
}

func GetStorageBucketIamPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]Asset, error) {
	return newStorageBucketIamAsset(d, config, expandIamPolicyBindings)
}

func GetStorageBucketIamBindingCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]Asset, error) {
	return newStorageBucketIamAsset(d, config, expandIamRoleBindings)
}

func GetStorageBucketIamMemberCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]Asset, error) {
	return newStorageBucketIamAsset(d, config, expandIamMemberBindings)
}

func MergeStorageBucketIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeStorageBucketIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeStorageBucketIamBindingDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAuthoritativeBindings)
}

func MergeStorageBucketIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func MergeStorageBucketIamMemberDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAdditiveBindings)
}

func newStorageBucketIamAsset(
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	expandBindings func(d tpgresource.TerraformResourceData) ([]IAMBinding, error),
) ([]Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := assetName(d, config, "//storage.googleapis.com/{{bucket}}")
	if err != nil {
		return []Asset{}, err
	}

	return []Asset{{
		Name: name,
		Type: StorageBucketIAMAssetType,
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchStorageBucketIamPolicy(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("bucket"); !ok {
		return Asset{}, ErrEmptyIdentityField
	}

	return fetchIamPolicy(
		StorageBucketIamUpdaterProducer,
		d,
		config,
		"//storage.googleapis.com/{{bucket}}",
		StorageBucketIAMAssetType,
	)
}
