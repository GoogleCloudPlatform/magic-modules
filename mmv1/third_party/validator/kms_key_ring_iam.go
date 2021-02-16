package google

import "fmt"

func GetKmsCryptoKeyIamPolicyCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newKmsCryptoKeyIamAsset(d, config, expandIamPolicyBindings)
}

func GetKmsCryptoKeyIamBindingCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newKmsCryptoKeyIamAsset(d, config, expandIamRoleBindings)
}

func GetKmsCryptoKeyIamMemberCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newKmsCryptoKeyIamAsset(d, config, expandIamMemberBindings)
}

func MergeKmsCryptoKeyIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeKmsCryptoKeyIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeKmsCryptoKeyIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func newKmsCryptoKeyIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) (Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	// Ideally we should use KmsCryptoKey_number, but since that is generated server-side,
	// we substitute KmsCryptoKey_id.
	name, err := assetName(d, config, "//cloudkms.googleapis.com/projects/{{project}}/locations/{{location}}/keyRings/{{key_ring}}/cryptoKeys/{{name}}")
	if err != nil {
		return Asset{}, err
	}

	return Asset{
		Name: name,
		Type: "cloudkms.googleapis.com/CryptoKey",
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}, nil
}
