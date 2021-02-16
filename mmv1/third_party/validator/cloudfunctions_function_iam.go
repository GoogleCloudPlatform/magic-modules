package google

import "fmt"

func GetCloudFunctionFunctionIamPolicyCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return GetCloudFunctionFunctionIamAsset(d, config, expandIamPolicyBindings)
}

func GetCloudFunctionFunctionIamBindingCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return GetCloudFunctionFunctionIamAsset(d, config, expandIamRoleBindings)
}

func GetCloudFunctionFunctionIamMemberCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return GetCloudFunctionFunctionIamAsset(d, config, expandIamMemberBindings)
}

func MergeCloudFunctionFunctionIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeCloudFunctionFunctionIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeCloudFunctionFunctionIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func newGetCloudFunctionFunctionIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) (Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	// Ideally we should use project_number, but since that is generated server-side,
	// we substitute project_id.
	name, err := assetName(d, config, "//cloudfunctions.googleapis.com/projects/{{project}}/locations/{{region}}/functions/{{name}}")
	if err != nil {
		return Asset{}, err
	}

	return Asset{
		Name: name,
		Type: "cloudfunctions.googleapis.com/CloudFunction",
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}, nil
}
