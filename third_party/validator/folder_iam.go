package google

import "fmt"

func GetFolderIamPolicyCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newFolderIamAsset(d, config, expandIamPolicyBindings)
}

func GetFolderIamBindingCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newFolderIamAsset(d, config, expandIamRoleBindings)
}

func GetFolderIamMemberCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newFolderIamAsset(d, config, expandIamMemberBindings)
}

func MergeFolderIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeFolderIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeFolderIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func newFolderIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) (Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	return Asset{
		// The "folder" argument is of the form "folders/12345"
		Name: fmt.Sprintf("//cloudresourcemanager.googleapis.com/%v", d.Get("folder").(string)),
		Type: "cloudresourcemanager.googleapis.com/Folder",
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}, nil
}
