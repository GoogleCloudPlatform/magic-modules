package pubsub

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

func ResourceConverterPubsubSubscriptionIamPolicy() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "pubsub.googleapis.com/Subscription",
		Convert:           GetPubsubSubscriptionIamPolicyCaiObject,
		MergeCreateUpdate: MergePubsubSubscriptionIamPolicy,
	}
}

func ResourceConverterPubsubSubscriptionIamBinding() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "pubsub.googleapis.com/Subscription",
		Convert:           GetPubsubSubscriptionIamBindingCaiObject,
		FetchFullResource: FetchPubsubSubscriptionIamPolicy,
		MergeCreateUpdate: MergePubsubSubscriptionIamBinding,
		MergeDelete:       MergePubsubSubscriptionIamBindingDelete,
	}
}

func ResourceConverterPubsubSubscriptionIamMember() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "pubsub.googleapis.com/Subscription",
		Convert:           GetPubsubSubscriptionIamMemberCaiObject,
		FetchFullResource: FetchPubsubSubscriptionIamPolicy,
		MergeCreateUpdate: MergePubsubSubscriptionIamMember,
		MergeDelete:       MergePubsubSubscriptionIamMemberDelete,
	}
}

func GetPubsubSubscriptionIamPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newPubsubSubscriptionIamAsset(d, config, cai.ExpandIamPolicyBindings)
}

func GetPubsubSubscriptionIamBindingCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newPubsubSubscriptionIamAsset(d, config, cai.ExpandIamRoleBindings)
}

func GetPubsubSubscriptionIamMemberCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newPubsubSubscriptionIamAsset(d, config, cai.ExpandIamMemberBindings)
}

func MergePubsubSubscriptionIamPolicy(existing, incoming cai.Asset) cai.Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergePubsubSubscriptionIamBinding(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeIamAssets(existing, incoming, cai.MergeAuthoritativeBindings)
}

func MergePubsubSubscriptionIamBindingDelete(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeDeleteIamAssets(existing, incoming, cai.MergeDeleteAuthoritativeBindings)
}

func MergePubsubSubscriptionIamMember(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeIamAssets(existing, incoming, cai.MergeAdditiveBindings)
}

func MergePubsubSubscriptionIamMemberDelete(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeDeleteIamAssets(existing, incoming, cai.MergeDeleteAdditiveBindings)
}

func newPubsubSubscriptionIamAsset(
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	expandBindings func(d tpgresource.TerraformResourceData) ([]cai.IAMBinding, error),
) ([]cai.Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []cai.Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := cai.AssetName(d, config, "//pubsub.googleapis.com/projects/{{project}}/subscriptions/{{subscription}}")
	if err != nil {
		return []cai.Asset{}, err
	}

	return []cai.Asset{{
		Name: name,
		Type: "pubsub.googleapis.com/Subscription",
		IAMPolicy: &cai.IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchPubsubSubscriptionIamPolicy(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (cai.Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("subscription"); !ok {
		return cai.Asset{}, cai.ErrEmptyIdentityField
	}

	return cai.FetchIamPolicy(
		NewPubsubSubscriptionIamUpdater,
		d,
		config,
		"//pubsub.googleapis.com/projects/{{project}}/subscriptions/{{subscription}}",
		"pubsub.googleapis.com/Subscription",
	)
}
