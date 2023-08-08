package spanner

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgiamresource"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
)

func ResourceConverterSpannerInstanceIamPolicy() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "spanner.googleapis.com/Instance",
		Convert:           GetSpannerInstanceIamPolicyCaiObject,
		MergeCreateUpdate: MergeSpannerInstanceIamPolicy,
	}
}

func ResourceConverterSpannerInstanceIamBinding() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "spanner.googleapis.com/Instance",
		Convert:           GetSpannerInstanceIamBindingCaiObject,
		FetchFullResource: FetchSpannerInstanceIamPolicy,
		MergeCreateUpdate: MergeSpannerInstanceIamBinding,
		MergeDelete:       MergeSpannerInstanceIamBindingDelete,
	}
}

func ResourceConverterSpannerInstanceIamMember() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "spanner.googleapis.com/Instance",
		Convert:           GetSpannerInstanceIamMemberCaiObject,
		FetchFullResource: FetchSpannerInstanceIamPolicy,
		MergeCreateUpdate: MergeSpannerInstanceIamMember,
		MergeDelete:       MergeSpannerInstanceIamMemberDelete,
	}
}

func GetSpannerInstanceIamPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newSpannerInstanceIamAsset(d, config, tpgiamresource.ExpandIamPolicyBindings)
}

func GetSpannerInstanceIamBindingCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newSpannerInstanceIamAsset(d, config, tpgiamresource.ExpandIamRoleBindings)
}

func GetSpannerInstanceIamMemberCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newSpannerInstanceIamAsset(d, config, tpgiamresource.ExpandIamMemberBindings)
}

func MergeSpannerInstanceIamPolicy(existing, incoming tpgresource.Asset) tpgresource.Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeSpannerInstanceIamBinding(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeIamAssets(existing, incoming, tpgiamresource.MergeAuthoritativeBindings)
}

func MergeSpannerInstanceIamBindingDelete(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeDeleteIamAssets(existing, incoming, tpgiamresource.MergeDeleteAuthoritativeBindings)
}

func MergeSpannerInstanceIamMember(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeIamAssets(existing, incoming, tpgiamresource.MergeAdditiveBindings)
}

func MergeSpannerInstanceIamMemberDelete(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeDeleteIamAssets(existing, incoming, tpgiamresource.MergeDeleteAdditiveBindings)
}

func newSpannerInstanceIamAsset(
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	expandBindings func(d tpgresource.TerraformResourceData) ([]tpgresource.IAMBinding, error),
) ([]tpgresource.Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []tpgresource.Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := tpgresource.AssetName(d, config, "//spanner.googleapis.com/projects/{{project}}/instances/{{instance}}")
	if err != nil {
		return []tpgresource.Asset{}, err
	}

	return []tpgresource.Asset{{
		Name: name,
		Type: "spanner.googleapis.com/Instance",
		IAMPolicy: &tpgresource.IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchSpannerInstanceIamPolicy(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgresource.Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("instance"); !ok {
		return tpgresource.Asset{}, tpgresource.ErrEmptyIdentityField
	}

	return tpgiamresource.FetchIamPolicy(
		NewSpannerInstanceIamUpdater,
		d,
		config,
		"//spanner.googleapis.com/projects/{{project}}/instances/{{instance}}",
		"spanner.googleapis.com/Instance",
	)
}
