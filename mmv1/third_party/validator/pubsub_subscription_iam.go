package google

import "fmt"

func GetPubSubSubscriptionIamPolicyCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return GetPubSubSubscriptionIamAsset(d, config, expandIamPolicyBindings)
}

func GetPubSubSubscriptionIamBindingCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return GetPubSubSubscriptionIamAsset(d, config, expandIamRoleBindings)
}

func GetPubSubSubscriptionIamMemberCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	return GetPubSubSubscriptionIamAsset(d, config, expandIamMemberBindings)
}

func MergePubSubSubscriptionIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergePubSubSubscriptionIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergePubSubSubscriptionIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func newGetPubSubSubscriptionIamAsset(
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
	name, err := assetName(d, config, "//pubsub.googleapis.com/projects/{{project}}/subscriptions/{{name}}")
	if err != nil {
		return Asset{}, err
	}

	return Asset{
		Name: name,
		Type: "pubsub.googleapis.com/Subscription",
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}, nil
}
