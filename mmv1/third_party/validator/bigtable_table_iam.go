package google

import "fmt"

func GetBigtableTableIamPolicyCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newBigtableTableIamAsset(d, config, expandIamPolicyBindings)
}

func GetBigtableTableIamBindingCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newBigtableTableIamAsset(d, config, expandIamRoleBindings)
}

func GetBigtableTableIamMemberCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newBigtableTableIamAsset(d, config, expandIamMemberBindings)
}

func MergeBigtableTableIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeBigtableTableIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeBigtableTableIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func newBigtableTableIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) (Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	// Ideally we should use BigtableTable_number, but since that is generated server-side,
	// we substitute BigtableTable_id.
	name, err := assetName(d, config, "//bigtable.googleapis.com/projects/{{project}}/instances/{{instance_name}}/tables/{{name}}")
	if err != nil {
		return Asset{}, err
	}

	return Asset{
		Name: name,
		Type: "bigtableadmin.googleapis.com/Table",
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}, nil
}
