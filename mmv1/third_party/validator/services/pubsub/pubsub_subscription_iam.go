package pubsub

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgiamresource"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
)

func ResourceConverterPubsubSubscriptionIamPolicy() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "pubsub.googleapis.com/Subscription",
		Convert:           GetPubsubSubscriptionIamPolicyCaiObject,
		MergeCreateUpdate: MergePubsubSubscriptionIamPolicy,
	}
}

func ResourceConverterPubsubSubscriptionIamBinding() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "pubsub.googleapis.com/Subscription",
		Convert:           GetPubsubSubscriptionIamBindingCaiObject,
		FetchFullResource: FetchPubsubSubscriptionIamPolicy,
		MergeCreateUpdate: MergePubsubSubscriptionIamBinding,
		MergeDelete:       MergePubsubSubscriptionIamBindingDelete,
	}
}

func ResourceConverterPubsubSubscriptionIamMember() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "pubsub.googleapis.com/Subscription",
		Convert:           GetPubsubSubscriptionIamMemberCaiObject,
		FetchFullResource: FetchPubsubSubscriptionIamPolicy,
		MergeCreateUpdate: MergePubsubSubscriptionIamMember,
		MergeDelete:       MergePubsubSubscriptionIamMemberDelete,
	}
}

func GetPubsubSubscriptionIamPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newPubsubSubscriptionIamAsset(d, config, tpgiamresource.ExpandIamPolicyBindings)
}

func GetPubsubSubscriptionIamBindingCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newPubsubSubscriptionIamAsset(d, config, tpgiamresource.ExpandIamRoleBindings)
}

func GetPubsubSubscriptionIamMemberCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newPubsubSubscriptionIamAsset(d, config, tpgiamresource.ExpandIamMemberBindings)
}

func MergePubsubSubscriptionIamPolicy(existing, incoming tpgresource.Asset) tpgresource.Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergePubsubSubscriptionIamBinding(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeIamAssets(existing, incoming, tpgiamresource.MergeAuthoritativeBindings)
}

func MergePubsubSubscriptionIamBindingDelete(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeDeleteIamAssets(existing, incoming, tpgiamresource.MergeDeleteAuthoritativeBindings)
}

func MergePubsubSubscriptionIamMember(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeIamAssets(existing, incoming, tpgiamresource.MergeAdditiveBindings)
}

func MergePubsubSubscriptionIamMemberDelete(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeDeleteIamAssets(existing, incoming, tpgiamresource.MergeDeleteAdditiveBindings)
}

func newPubsubSubscriptionIamAsset(
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	expandBindings func(d tpgresource.TerraformResourceData) ([]tpgresource.IAMBinding, error),
) ([]tpgresource.Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []tpgresource.Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := tpgresource.AssetName(d, config, "//pubsub.googleapis.com/projects/{{project}}/subscriptions/{{subscription}}")
	if err != nil {
		return []tpgresource.Asset{}, err
	}

	return []tpgresource.Asset{{
		Name: name,
		Type: "pubsub.googleapis.com/Subscription",
		IAMPolicy: &tpgresource.IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchPubsubSubscriptionIamPolicy(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgresource.Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("subscription"); !ok {
		return tpgresource.Asset{}, tpgresource.ErrEmptyIdentityField
	}

	return tpgiamresource.FetchIamPolicy(
		NewPubsubSubscriptionIamUpdater,
		d,
		config,
		"//pubsub.googleapis.com/projects/{{project}}/subscriptions/{{subscription}}",
		"pubsub.googleapis.com/Subscription",
	)
}
