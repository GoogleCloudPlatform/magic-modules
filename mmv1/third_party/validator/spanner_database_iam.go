package google

import "fmt"

func resourceConverterSpannerDatabaseIamPolicy() ResourceConverter {
	return ResourceConverter{
		AssetType:         "spanner.googleapis.com/Database",
		Convert:           GetSpannerDatabaseIamPolicyCaiObject,
		MergeCreateUpdate: MergeSpannerDatabaseIamPolicy,
	}
}

func resourceConverterSpannerDatabaseIamBinding() ResourceConverter {
	return ResourceConverter{
		AssetType:         "spanner.googleapis.com/Database",
		Convert:           GetSpannerDatabaseIamBindingCaiObject,
		FetchFullResource: FetchSpannerDatabaseIamPolicy,
		MergeCreateUpdate: MergeSpannerDatabaseIamBinding,
		MergeDelete:       MergeSpannerDatabaseIamBindingDelete,
	}
}

func resourceConverterSpannerDatabaseIamMember() ResourceConverter {
	return ResourceConverter{
		AssetType:         "spanner.googleapis.com/Database",
		Convert:           GetSpannerDatabaseIamMemberCaiObject,
		FetchFullResource: FetchSpannerDatabaseIamPolicy,
		MergeCreateUpdate: MergeSpannerDatabaseIamMember,
		MergeDelete:       MergeSpannerDatabaseIamMemberDelete,
	}
}

func GetSpannerDatabaseIamPolicyCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newSpannerDatabaseIamAsset(d, config, expandIamPolicyBindings)
}

func GetSpannerDatabaseIamBindingCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newSpannerDatabaseIamAsset(d, config, expandIamRoleBindings)
}

func GetSpannerDatabaseIamMemberCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newSpannerDatabaseIamAsset(d, config, expandIamMemberBindings)
}

func MergeSpannerDatabaseIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeSpannerDatabaseIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeSpannerDatabaseIamBindingDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAuthoritativeBindings)
}

func MergeSpannerDatabaseIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func MergeSpannerDatabaseIamMemberDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAdditiveBindings)
}

func newSpannerDatabaseIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) ([]Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := assetName(d, config, "//spanner.googleapis.com/projects/{{project}}/instances/{{instance}}/databases/{{database}}")
	if err != nil {
		return []Asset{}, err
	}

	return []Asset{{
		Name: name,
		Type: "spanner.googleapis.com/Database",
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchSpannerDatabaseIamPolicy(d TerraformResourceData, config *Config) (Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("instance"); !ok {
		return Asset{}, ErrEmptyIdentityField
	}

	if _, ok := d.GetOk("database"); !ok {
		return Asset{}, ErrEmptyIdentityField
	}

	return fetchIamPolicy(
		NewSpannerDatabaseIamUpdater,
		d,
		config,
		"//spanner.googleapis.com/projects/{{project}}/instances/{{instance}}/databases/{{database}}",
		"spanner.googleapis.com/Database",
	)
}
