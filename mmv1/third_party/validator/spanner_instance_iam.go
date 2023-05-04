package google

import (
	"fmt"

	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
)

func resourceConverterSpannerInstanceIamPolicy() ResourceConverter {
	return ResourceConverter{
		AssetType:         "spanner.googleapis.com/Instance",
		Convert:           GetSpannerInstanceIamPolicyCaiObject,
		MergeCreateUpdate: MergeSpannerInstanceIamPolicy,
	}
}

func resourceConverterSpannerInstanceIamBinding() ResourceConverter {
	return ResourceConverter{
		AssetType:         "spanner.googleapis.com/Instance",
		Convert:           GetSpannerInstanceIamBindingCaiObject,
		FetchFullResource: FetchSpannerInstanceIamPolicy,
		MergeCreateUpdate: MergeSpannerInstanceIamBinding,
		MergeDelete:       MergeSpannerInstanceIamBindingDelete,
	}
}

func resourceConverterSpannerInstanceIamMember() ResourceConverter {
	return ResourceConverter{
		AssetType:         "spanner.googleapis.com/Instance",
		Convert:           GetSpannerInstanceIamMemberCaiObject,
		FetchFullResource: FetchSpannerInstanceIamPolicy,
		MergeCreateUpdate: MergeSpannerInstanceIamMember,
		MergeDelete:       MergeSpannerInstanceIamMemberDelete,
	}
}

func GetSpannerInstanceIamPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]Asset, error) {
	return newSpannerInstanceIamAsset(d, config, expandIamPolicyBindings)
}

func GetSpannerInstanceIamBindingCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]Asset, error) {
	return newSpannerInstanceIamAsset(d, config, expandIamRoleBindings)
}

func GetSpannerInstanceIamMemberCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]Asset, error) {
	return newSpannerInstanceIamAsset(d, config, expandIamMemberBindings)
}

func MergeSpannerInstanceIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeSpannerInstanceIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeSpannerInstanceIamBindingDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAuthoritativeBindings)
}

func MergeSpannerInstanceIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func MergeSpannerInstanceIamMemberDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAdditiveBindings)
}

func newSpannerInstanceIamAsset(
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	expandBindings func(d tpgresource.TerraformResourceData) ([]IAMBinding, error),
) ([]Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := assetName(d, config, "//spanner.googleapis.com/projects/{{project}}/instances/{{instance}}")
	if err != nil {
		return []Asset{}, err
	}

	return []Asset{{
		Name: name,
		Type: "spanner.googleapis.com/Instance",
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchSpannerInstanceIamPolicy(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("instance"); !ok {
		return Asset{}, ErrEmptyIdentityField
	}

	return fetchIamPolicy(
		NewSpannerInstanceIamUpdater,
		d,
		config,
		"//spanner.googleapis.com/projects/{{project}}/instances/{{instance}}",
		"spanner.googleapis.com/Instance",
	)
}
