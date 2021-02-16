package google

import "fmt"

func GetBigQueryTableIamPolicyCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newBigQueryTableIamAsset(d, config, expandIamPolicyBindings)
}

func GetBigQueryTableIamBindingCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newBigQueryTableIamAsset(d, config, expandIamRoleBindings)
}

func GetBigQueryTableIamMemberCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newBigQueryTableIamAsset(d, config, expandIamMemberBindings)
}

func MergeBigQueryTableIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeBigQueryTableIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeBigQueryTableIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func newBigQueryTableIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) (Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	// Ideally we should use BigQueryTable_number, but since that is generated server-side,
	// we substitute BigQueryTable_id.
	name, err := assetName(d, config, "//bigquery.googleapis.com/projects/{{.Provider.project}}/datasets/{{dataset_id}}/tables/{{table_id}}")
	if err != nil {
		return Asset{}, err
	}

	return Asset{
		Name: name,
		Type: "bigquery.googleapis.com/Table",
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}, nil
}
