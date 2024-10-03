package storage

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

// Provide a separate asset type constant so we don't have to worry about name conflicts between IAM and non-IAM converter files
const StorageBucketIAMAssetType string = "storage.googleapis.com/Bucket"

func ResourceConverterStorageBucketIamPolicy() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         StorageBucketIAMAssetType,
		Convert:           GetStorageBucketIamPolicyCaiObject,
		MergeCreateUpdate: MergeStorageBucketIamPolicy,
	}
}

func ResourceConverterStorageBucketIamBinding() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         StorageBucketIAMAssetType,
		Convert:           GetStorageBucketIamBindingCaiObject,
		FetchFullResource: FetchStorageBucketIamPolicy,
		MergeCreateUpdate: MergeStorageBucketIamBinding,
		MergeDelete:       MergeStorageBucketIamBindingDelete,
	}
}

func ResourceConverterStorageBucketIamMember() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         StorageBucketIAMAssetType,
		Convert:           GetStorageBucketIamMemberCaiObject,
		FetchFullResource: FetchStorageBucketIamPolicy,
		MergeCreateUpdate: MergeStorageBucketIamMember,
		MergeDelete:       MergeStorageBucketIamMemberDelete,
	}
}

func GetStorageBucketIamPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newStorageBucketIamAsset(d, config, cai.ExpandIamPolicyBindings)
}

func GetStorageBucketIamBindingCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newStorageBucketIamAsset(d, config, cai.ExpandIamRoleBindings)
}

func GetStorageBucketIamMemberCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newStorageBucketIamAsset(d, config, cai.ExpandIamMemberBindings)
}

func MergeStorageBucketIamPolicy(existing, incoming cai.Asset) cai.Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeStorageBucketIamBinding(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeIamAssets(existing, incoming, cai.MergeAuthoritativeBindings)
}

func MergeStorageBucketIamBindingDelete(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeDeleteIamAssets(existing, incoming, cai.MergeDeleteAuthoritativeBindings)
}

func MergeStorageBucketIamMember(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeIamAssets(existing, incoming, cai.MergeAdditiveBindings)
}

func MergeStorageBucketIamMemberDelete(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeDeleteIamAssets(existing, incoming, cai.MergeDeleteAdditiveBindings)
}

func newStorageBucketIamAsset(
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	expandBindings func(d tpgresource.TerraformResourceData) ([]cai.IAMBinding, error),
) ([]cai.Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []cai.Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := cai.AssetName(d, config, "//storage.googleapis.com/{{bucket}}")
	if err != nil {
		return []cai.Asset{}, err
	}

	return []cai.Asset{{
		Name: name,
		Type: StorageBucketIAMAssetType,
		IAMPolicy: &cai.IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchStorageBucketIamPolicy(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (cai.Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("bucket"); !ok {
		return cai.Asset{}, cai.ErrEmptyIdentityField
	}

	return cai.FetchIamPolicy(
		StorageBucketIamUpdaterProducer,
		d,
		config,
		"//storage.googleapis.com/{{bucket}}",
		StorageBucketIAMAssetType,
	)
}
