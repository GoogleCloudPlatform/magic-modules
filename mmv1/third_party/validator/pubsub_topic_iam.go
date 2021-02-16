package google

import "fmt"

func GetPubSubTopicIamPolicyCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newPubSubTopicIamAsset(d, config, expandIamPolicyBindings)
}

func GetPubSubTopicIamBindingCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newPubSubTopicIamAsset(d, config, expandIamRoleBindings)
}

func GetPubSubTopicIamMemberCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return newPubSubTopicIamAsset(d, config, expandIamMemberBindings)
}

func MergePubSubTopicIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergePubSubTopicIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergePubSubTopicIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func newPubSubTopicIamAsset(
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
	name, err := assetName(d, config, "//pubsub.googleapis.com/projects/{{project}}/topics/{{topic}}")
	if err != nil {
		return Asset{}, err
	}

	return Asset{
		Name: name,
		Type: "pubsub.googleapis.com/Topic",
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}, nil
}
