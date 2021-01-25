package google

import "fmt"

func GetBucketIamPolicyCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newBucketIamAsset(d, config, expandIamPolicyBindings)
}

func GetBucketIamBindingCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newBucketIamAsset(d, config, expandIamRoleBindings)
}

func GetBucketIamMemberCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newBucketIamAsset(d, config, expandIamMemberBindings)
}

func MergeBucketIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeBucketIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeBucketIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func newBucketIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) (Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := assetName(d, config, "//storage.googleapis.com/{{name}}")
	if err != nil {
		return Asset{}, err
	}

	return Asset{
		Name: name,
		Type: "storage.googleapis.com/Bucket",
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}, nil
}
