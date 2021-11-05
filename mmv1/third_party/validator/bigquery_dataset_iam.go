package google

import "fmt"

func resourceConverterBigQueryDatasetIamPolicy() ResourceConverter {
	return ResourceConverter{
		AssetType:         "bigquery.googleapis.com/Dataset",
		Convert:           GetBigQueryDatasetIamPolicyCaiObject,
		MergeCreateUpdate: MergeBigQueryDatasetIamPolicy,
	}
}

func resourceConverterBigQueryDatasetIamBinding() ResourceConverter {
	return ResourceConverter{
		AssetType:         "bigquery.googleapis.com/Dataset",
		Convert:           GetBigQueryDatasetIamBindingCaiObject,
		FetchFullResource: FetchBigQueryDatasetIamPolicy,
		MergeCreateUpdate: MergeBigQueryDatasetIamBinding,
		MergeDelete:       MergeBigQueryDatasetIamBindingDelete,
	}
}

func resourceConverterBigQueryDatasetIamMember() ResourceConverter {
	return ResourceConverter{
		AssetType:         "bigquery.googleapis.com/Dataset",
		Convert:           GetBigQueryDatasetIamMemberCaiObject,
		FetchFullResource: FetchBigQueryDatasetIamPolicy,
		MergeCreateUpdate: MergeBigQueryDatasetIamMember,
		MergeDelete:       MergeBigQueryDatasetIamMemberDelete,
	}
}

func GetBigQueryDatasetIamPolicyCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newBigQueryDatasetIamAsset(d, config, expandIamPolicyBindings)
}

func GetBigQueryDatasetIamBindingCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newBigQueryDatasetIamAsset(d, config, expandIamRoleBindings)
}

func GetBigQueryDatasetIamMemberCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newBigQueryDatasetIamAsset(d, config, expandIamMemberBindings)
}

func MergeBigQueryDatasetIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeBigQueryDatasetIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeBigQueryDatasetIamBindingDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAuthoritativeBindings)
}

func MergeBigQueryDatasetIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func MergeBigQueryDatasetIamMemberDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAdditiveBindings)
}

func newBigQueryDatasetIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) ([]Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := assetName(d, config, "//bigquery.googleapis.com/datasets/{{dataset_id}}")
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

func FetchBigQueryDatasetIamPolicy(d TerraformResourceData, config *Config) (Asset, error) {
	return fetchIamPolicy(
		NewBigQueryDatasetIamUpdater,
		d,
		config,
		"//bigquery.googleapis.com/datasets/{{dataset_id}}",
		"bigquery.googleapis.com/Dataset",
	)
}
