package google

import "fmt"

func GetBucketIamPolicyCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newBucketIamAsset(d, config, expandIamPolicyBindings)
}

func GetBucketIamBindingCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newBucketIamAsset(d, config, expandIamRoleBindings)
}

func GetBucketIamMemberCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newBucketIamAsset(d, config, expandIamMemberBindings)
}

func MergeBucketIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeBucketIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeBucketIamBindingDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAuthoritativeBindings)
}

func MergeBucketIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func MergeBucketIamMemberDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAdditiveBindings)
}

func newBucketIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
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
		Type: "storage.googleapis.com/Bucket",
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchBucketIamPolicy(d TerraformResourceData, config *Config) (Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("{{bucket}}"); !ok {
		return Asset{}, ErrEmptyIdentityField
	}

	return fetchIamPolicy(
		StorageBucketIamUpdaterProducer,
		d,
		config,
		"//storage.googleapis.com/{{bucket}}",
		"storage.googleapis.com/Bucket",
	)
}
