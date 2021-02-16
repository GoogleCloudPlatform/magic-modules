package google

import "fmt"

func GetBigtableInstanceIamPolicyCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newBigtableInstanceIamAsset(d, config, expandIamPolicyBindings)
}

func GetBigtableInstanceIamBindingCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newBigtableInstanceIamAsset(d, config, expandIamRoleBindings)
}

func GetBigtableInstanceIamMemberCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newBigtableInstanceIamAsset(d, config, expandIamMemberBindings)
}

func MergeBigtableInstanceIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeBigtableInstanceIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeBigtableInstanceIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func newBigtableInstanceIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) (Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	// Ideally we should use BigtableInstance_number, but since that is generated server-side,
	// we substitute BigtableInstance_id.
	name, err := assetName(d, config, "//bigtable.googleapis.com/projects/{{project}}/instances/{{name}}")
	if err != nil {
		return Asset{}, err
	}

	return Asset{
		Name: name,
		Type: "bigtableadmin.googleapis.com/Instance",
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}, nil
}
