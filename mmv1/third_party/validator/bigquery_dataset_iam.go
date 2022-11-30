package google

import "fmt"

func resourceConverterBigqueryDatasetIamPolicy() ResourceConverter {
	return ResourceConverter{
		AssetType:         "bigquery.googleapis.com/Dataset",
		Convert:           GetBigqueryDatasetIamPolicyCaiObject,
		MergeCreateUpdate: MergeBigqueryDatasetIamPolicy,
	}
}

func resourceConverterBigqueryDatasetIamBinding() ResourceConverter {
	return ResourceConverter{
		AssetType:         "bigquery.googleapis.com/Dataset",
		Convert:           GetBigqueryDatasetIamBindingCaiObject,
		FetchFullResource: FetchBigqueryDatasetIamPolicy,
		MergeCreateUpdate: MergeBigqueryDatasetIamBinding,
		MergeDelete:       MergeBigqueryDatasetIamBindingDelete,
	}
}

func resourceConverterBigqueryDatasetIamMember() ResourceConverter {
	return ResourceConverter{
		AssetType:         "bigquery.googleapis.com/Dataset",
		Convert:           GetBigqueryDatasetIamMemberCaiObject,
		FetchFullResource: FetchBigqueryDatasetIamPolicy,
		MergeCreateUpdate: MergeBigqueryDatasetIamMember,
		MergeDelete:       MergeBigqueryDatasetIamMemberDelete,
	}
}

func GetBigqueryDatasetIamPolicyCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newBigqueryDatasetIamAsset(d, config, expandIamPolicyBindings)
}

func GetBigqueryDatasetIamBindingCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newBigqueryDatasetIamAsset(d, config, expandIamRoleBindings)
}

func GetBigqueryDatasetIamMemberCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newBigqueryDatasetIamAsset(d, config, expandIamMemberBindings)
}

func MergeBigqueryDatasetIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeBigqueryDatasetIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeBigqueryDatasetIamBindingDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAuthoritativeBindings)
}

func MergeBigqueryDatasetIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func MergeBigqueryDatasetIamMemberDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAdditiveBindings)
}

func newBigqueryDatasetIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) ([]Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := assetName(d, config, "//bigquery.googleapis.com/projects/{{project}}/datasets/{{dataset_id}}")
	if err != nil {
		return []Asset{}, err
	}

	return []Asset{{
		Name: name,
		Type: "bigquery.googleapis.com/Dataset",
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchBigqueryDatasetIamPolicy(d TerraformResourceData, config *Config) (Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("dataset_id"); !ok {
		return Asset{}, ErrEmptyIdentityField
	}

	return fetchIamPolicy(
		NewBigqueryDatasetIamUpdater,
		d,
		config,
		"//bigquery.googleapis.com/projects/{{project}}/datasets/{{dataset_id}}",
		"bigquery.googleapis.com/Dataset",
	)
}
