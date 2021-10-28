package google

import "fmt"

func resourceConverterKmsKeyRingIamPolicy() ResourceConverter {
	return ResourceConverter{
		AssetType:         "cloudkms.googleapis.com/KeyRing",
		Convert:           GetKmsKeyRingIamPolicyCaiObject,
		MergeCreateUpdate: MergeKmsKeyRingIamPolicy,
	}
}

func resourceConverterKmsKeyRingIamBinding() ResourceConverter {
	return ResourceConverter{
		AssetType:         "cloudkms.googleapis.com/KeyRing",
		Convert:           GetKmsKeyRingIamBindingCaiObject,
		FetchFullResource: FetchKmsKeyRingIamPolicy,
		MergeCreateUpdate: MergeKmsKeyRingIamBinding,
		MergeDelete:       MergeKmsKeyRingIamBindingDelete,
	}
}

func resourceConverterKmsKeyRingIamMember() ResourceConverter {
	return ResourceConverter{
		AssetType:         "cloudkms.googleapis.com/KeyRing",
		Convert:           GetKmsKeyRingIamMemberCaiObject,
		FetchFullResource: FetchKmsKeyRingIamPolicy,
		MergeCreateUpdate: MergeKmsKeyRingIamMember,
		MergeDelete:       MergeKmsKeyRingIamMemberDelete,
	}
}

func GetKmsKeyRingIamPolicyCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newKmsKeyRingIamAsset(d, config, expandIamPolicyBindings)
}

func GetKmsKeyRingIamBindingCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newKmsKeyRingIamAsset(d, config, expandIamRoleBindings)
}

func GetKmsKeyRingIamMemberCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newKmsKeyRingIamAsset(d, config, expandIamMemberBindings)
}

func MergeKmsKeyRingIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeKmsKeyRingIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeKmsKeyRingIamBindingDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAuthoritativeBindings)
}

func MergeKmsKeyRingIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func MergeKmsKeyRingIamMemberDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAdditiveBindings)
}

func newKmsKeyRingIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) ([]Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := assetName(d, config, "//cloudkms.googleapis.com/{{key_ring_id}}")
	if err != nil {
		return []Asset{}, err
	}

	return []Asset{{
		Name: name,
		Type: "cloudkms.googleapis.com/KeyRing",
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchKmsKeyRingIamPolicy(d TerraformResourceData, config *Config) (Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("{{key_ring_id}}"); !ok {
		return Asset{}, ErrEmptyIdentityField
	}

	// We use key_ring_id in the asset name template to be consistent with newKmsKeyRingIamAsset.
	return fetchIamPolicy(
		NewKmsKeyRingIamUpdater,
		d,
		config,
		"//cloudkms.googleapis.com/{{key_ring_id}}", // asset name
		"cloudkms.googleapis.com/KeyRing",           // asset type
	)
}
