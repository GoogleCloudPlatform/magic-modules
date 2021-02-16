package google

import "fmt"

func GetKmsKeyRingIamPolicyCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newKmsKeyRingIamAsset(d, config, expandIamPolicyBindings)
}

func GetKmsKeyRingIamBindingCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newKmsKeyRingIamAsset(d, config, expandIamRoleBindings)
}

func GetKmsKeyRingIamMemberCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newKmsKeyRingIamAsset(d, config, expandIamMemberBindings)
}

func MergeKmsKeyRingIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeKmsKeyRingIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeKmsKeyRingIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func newKmsKeyRingIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) (Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	// Ideally we should use KmsKeyRing_number, but since that is generated server-side,
	// we substitute KmsKeyRing_id.
	name, err := assetName(d, config, "//cloudkms.googleapis.com/projects/{{project}}/locations/{{location}}/keyRings/{{name}}")
	if err != nil {
		return Asset{}, err
	}

	return Asset{
		Name: name,
		Type: "cloudkms.googleapis.com/KeyRing",
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}, nil
}
