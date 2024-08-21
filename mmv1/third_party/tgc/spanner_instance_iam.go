package spanner

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

func ResourceConverterSpannerInstanceIamPolicy() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "spanner.googleapis.com/Instance",
		Convert:           GetSpannerInstanceIamPolicyCaiObject,
		MergeCreateUpdate: MergeSpannerInstanceIamPolicy,
	}
}

func ResourceConverterSpannerInstanceIamBinding() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "spanner.googleapis.com/Instance",
		Convert:           GetSpannerInstanceIamBindingCaiObject,
		FetchFullResource: FetchSpannerInstanceIamPolicy,
		MergeCreateUpdate: MergeSpannerInstanceIamBinding,
		MergeDelete:       MergeSpannerInstanceIamBindingDelete,
	}
}

func ResourceConverterSpannerInstanceIamMember() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "spanner.googleapis.com/Instance",
		Convert:           GetSpannerInstanceIamMemberCaiObject,
		FetchFullResource: FetchSpannerInstanceIamPolicy,
		MergeCreateUpdate: MergeSpannerInstanceIamMember,
		MergeDelete:       MergeSpannerInstanceIamMemberDelete,
	}
}

func GetSpannerInstanceIamPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newSpannerInstanceIamAsset(d, config, cai.ExpandIamPolicyBindings)
}

func GetSpannerInstanceIamBindingCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newSpannerInstanceIamAsset(d, config, cai.ExpandIamRoleBindings)
}

func GetSpannerInstanceIamMemberCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newSpannerInstanceIamAsset(d, config, cai.ExpandIamMemberBindings)
}

func MergeSpannerInstanceIamPolicy(existing, incoming cai.Asset) cai.Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeSpannerInstanceIamBinding(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeIamAssets(existing, incoming, cai.MergeAuthoritativeBindings)
}

func MergeSpannerInstanceIamBindingDelete(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeDeleteIamAssets(existing, incoming, cai.MergeDeleteAuthoritativeBindings)
}

func MergeSpannerInstanceIamMember(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeIamAssets(existing, incoming, cai.MergeAdditiveBindings)
}

func MergeSpannerInstanceIamMemberDelete(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeDeleteIamAssets(existing, incoming, cai.MergeDeleteAdditiveBindings)
}

func newSpannerInstanceIamAsset(
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	expandBindings func(d tpgresource.TerraformResourceData) ([]cai.IAMBinding, error),
) ([]cai.Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []cai.Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := cai.AssetName(d, config, "//spanner.googleapis.com/projects/{{project}}/instances/{{instance}}")
	if err != nil {
		return []cai.Asset{}, err
	}

	return []cai.Asset{{
		Name: name,
		Type: "spanner.googleapis.com/Instance",
		IAMPolicy: &cai.IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchSpannerInstanceIamPolicy(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (cai.Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("instance"); !ok {
		return cai.Asset{}, cai.ErrEmptyIdentityField
	}

	return cai.FetchIamPolicy(
		NewSpannerInstanceIamUpdater,
		d,
		config,
		"//spanner.googleapis.com/projects/{{project}}/instances/{{instance}}",
		"spanner.googleapis.com/Instance",
	)
}
