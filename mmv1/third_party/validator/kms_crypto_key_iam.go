package google

import "fmt"

func GetKmsCryptoKeyIamPolicyCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newKmsCryptoKeyIamAsset(d, config, expandIamPolicyBindings)
}

func GetKmsCryptoKeyIamBindingCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newKmsCryptoKeyIamAsset(d, config, expandIamRoleBindings)
}

func GetKmsCryptoKeyIamMemberCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newKmsCryptoKeyIamAsset(d, config, expandIamMemberBindings)
}

func MergeKmsCryptoKeyIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeKmsCryptoKeyIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeKmsCryptoKeyIamBindingDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAuthoritativeBindings)
}

func MergeKmsCryptoKeyIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func MergeKmsCryptoKeyIamMemberDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAdditiveBindings)
}

func newKmsCryptoKeyIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) ([]Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := assetName(d, config, "//cloudkms.googleapis.com/{{crypto_key_id}}")
	if err != nil {
		return []Asset{}, err
	}

	return []Asset{{
		Name: name,
		Type: "cloudkms.googleapis.com/CryptoKey",
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchKmsCryptoKeyIamPolicy(d TerraformResourceData, config *Config) (Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("{{crypto_key_id}}"); !ok {
		return Asset{}, ErrEmptyIdentityField
	}

	// We use crypto_key_id in the asset name template to be consistent with newKmsCryptoKeyIamAsset.
	return fetchIamPolicy(
		NewKmsCryptoKeyIamUpdater,
		d,
		config,
		"//cloudkms.googleapis.com/{{crypto_key_id}}", // asset name
		"cloudkms.googleapis.com/CryptoKey",           // asset type
	)
}
