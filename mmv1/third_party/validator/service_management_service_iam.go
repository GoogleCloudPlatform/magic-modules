package google

import "fmt"

// Provide a separate asset type constant so we don't have to worry about name conflicts between IAM and non-IAM converter files
const ServiceManagementServiceIAMAssetType string = "servicemanagement.googleapis.com/Service"

func resourceConverterServiceManagementServiceIamPolicy() ResourceConverter {
	return ResourceConverter{
		AssetType:         ServiceManagementServiceIAMAssetType,
		Convert:           GetServiceManagementServiceIamPolicyCaiObject,
		MergeCreateUpdate: MergeServiceManagementServiceIamPolicy,
	}
}

func resourceConverterServiceManagementServiceIamBinding() ResourceConverter {
	return ResourceConverter{
		AssetType:         ServiceManagementServiceIAMAssetType,
		Convert:           GetServiceManagementServiceIamBindingCaiObject,
		FetchFullResource: FetchServiceManagementServiceIamPolicy,
		MergeCreateUpdate: MergeServiceManagementServiceIamBinding,
		MergeDelete:       MergeServiceManagementServiceIamBindingDelete,
	}
}

func resourceConverterServiceManagementServiceIamMember() ResourceConverter {
	return ResourceConverter{
		AssetType:         ServiceManagementServiceIAMAssetType,
		Convert:           GetServiceManagementServiceIamMemberCaiObject,
		FetchFullResource: FetchServiceManagementServiceIamPolicy,
		MergeCreateUpdate: MergeServiceManagementServiceIamMember,
		MergeDelete:       MergeServiceManagementServiceIamMemberDelete,
	}
}

func GetServiceManagementServiceIamPolicyCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newServiceManagementServiceIamAsset(d, config, expandIamPolicyBindings)
}

func GetServiceManagementServiceIamBindingCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newServiceManagementServiceIamAsset(d, config, expandIamRoleBindings)
}

func GetServiceManagementServiceIamMemberCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newServiceManagementServiceIamAsset(d, config, expandIamMemberBindings)
}

func MergeServiceManagementServiceIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeServiceManagementServiceIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeServiceManagementServiceIamBindingDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAuthoritativeBindings)
}

func MergeServiceManagementServiceIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func MergeServiceManagementServiceIamMemberDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAdditiveBindings)
}

func newServiceManagementServiceIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) ([]Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := assetName(d, config, "//servicemanagement.googleapis.com/services/{{service_name}}")
	if err != nil {
		return []Asset{}, err
	}

	return []Asset{{
		Name: name,
		Type: ServiceManagementServiceIAMAssetType,
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchServiceManagementServiceIamPolicy(d TerraformResourceData, config *Config) (Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("{{service_name}}"); !ok {
		return Asset{}, ErrEmptyIdentityField
	}

	return fetchIamPolicy(
		ServiceManagementServiceIamUpdaterProducer,
		d,
		config,
		"//servicemanagement.googleapis.com/services/{{service_name}}",
		ServiceManagementServiceIAMAssetType,
	)
}
