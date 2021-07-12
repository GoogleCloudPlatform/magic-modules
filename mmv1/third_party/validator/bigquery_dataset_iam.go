package google

import "fmt"

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

	// Ideally we should use BigqueryDataset_number, but since that is generated server-side,
	// we substitute BigqueryDataset_id.
	name, err := assetName(d, config, "//bigquery.googleapis.com/{{dataset_id}}")
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
	// We use BigqueryDataset_id in the asset name template to be consistent with newBigqueryDatasetIamAsset.
	return fetchIamPolicy(
		NewBigqueryDatasetIamUpdater,
		d,
		config,
		"//bigquery.googleapis.com/{{dataset_id}}", //asset name
		"bigquery.googleapis.com/Dataset",          //asset type
	)
}
