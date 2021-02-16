package google

import "fmt"

func GetBigQueryDatasetIamPolicyCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newBigQueryDatasetIamAsset(d, config, expandIamPolicyBindings)
}

func GetBigQueryDatasetIamBindingCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newBigQueryDatasetIamAsset(d, config, expandIamRoleBindings)
}

func GetBigQueryDatasetIamMemberCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newBigQueryDatasetIamAsset(d, config, expandIamMemberBindings)
}

func MergeBigQueryDatasetIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeBigQueryDatasetIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeBigQueryDatasetIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func newBigQueryDatasetIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) (Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	// Ideally we should use BigQueryDataset_number, but since that is generated server-side,
	// we substitute BigQueryDataset_id.
	name, err := assetName(d, config, "//bigquery.googleapis.com/projects/{{.Provider.project}}/datasets/{{dataset_id}}")
	if err != nil {
		return Asset{}, err
	}

	return Asset{
		Name: name,
		Type: "bigquery.googleapis.com/Dataset",
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}, nil
}
