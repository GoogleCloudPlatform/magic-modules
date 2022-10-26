package google

import "fmt"

func resourceConverterPubsubSubscriptionIamPolicy() ResourceConverter {
	return ResourceConverter{
		AssetType:         "pubsub.googleapis.com/Subscription",
		Convert:           GetPubsubSubscriptionIamPolicyCaiObject,
		MergeCreateUpdate: MergePubsubSubscriptionIamPolicy,
	}
}

func resourceConverterPubsubSubscriptionIamBinding() ResourceConverter {
	return ResourceConverter{
		AssetType:         "pubsub.googleapis.com/Subscription",
		Convert:           GetPubsubSubscriptionIamBindingCaiObject,
		FetchFullResource: FetchPubsubSubscriptionIamPolicy,
		MergeCreateUpdate: MergePubsubSubscriptionIamBinding,
		MergeDelete:       MergePubsubSubscriptionIamBindingDelete,
	}
}

func resourceConverterPubsubSubscriptionIamMember() ResourceConverter {
	return ResourceConverter{
		AssetType:         "pubsub.googleapis.com/Subscription",
		Convert:           GetPubsubSubscriptionIamMemberCaiObject,
		FetchFullResource: FetchPubsubSubscriptionIamPolicy,
		MergeCreateUpdate: MergePubsubSubscriptionIamMember,
		MergeDelete:       MergePubsubSubscriptionIamMemberDelete,
	}
}

func GetPubsubSubscriptionIamPolicyCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newPubsubSubscriptionIamAsset(d, config, expandIamPolicyBindings)
}

func GetPubsubSubscriptionIamBindingCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newPubsubSubscriptionIamAsset(d, config, expandIamRoleBindings)
}

func GetPubsubSubscriptionIamMemberCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newPubsubSubscriptionIamAsset(d, config, expandIamMemberBindings)
}

func MergePubsubSubscriptionIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergePubsubSubscriptionIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergePubsubSubscriptionIamBindingDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAuthoritativeBindings)
}

func MergePubsubSubscriptionIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func MergePubsubSubscriptionIamMemberDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAdditiveBindings)
}

func newPubsubSubscriptionIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) ([]Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := assetName(d, config, "//pubsub.googleapis.com/projects/{{project}}/subscriptions/{{subscription}}")
	if err != nil {
		return []Asset{}, err
	}

	return []Asset{{
		Name: name,
		Type: "pubsub.googleapis.com/Subscription",
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchPubsubSubscriptionIamPolicy(d TerraformResourceData, config *Config) (Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("subscription"); !ok {
		return Asset{}, ErrEmptyIdentityField
	}

	return fetchIamPolicy(
		NewPubsubSubscriptionIamUpdater,
		d,
		config,
		"//pubsub.googleapis.com/projects/{{project}}/subscriptions/{{subscription}}",
		"pubsub.googleapis.com/Subscription",
	)
}
