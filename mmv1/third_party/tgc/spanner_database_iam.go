package spanner

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

func ResourceConverterSpannerDatabaseIamPolicy() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "spanner.googleapis.com/Database",
		Convert:           GetSpannerDatabaseIamPolicyCaiObject,
		MergeCreateUpdate: MergeSpannerDatabaseIamPolicy,
	}
}

func ResourceConverterSpannerDatabaseIamBinding() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "spanner.googleapis.com/Database",
		Convert:           GetSpannerDatabaseIamBindingCaiObject,
		FetchFullResource: FetchSpannerDatabaseIamPolicy,
		MergeCreateUpdate: MergeSpannerDatabaseIamBinding,
		MergeDelete:       MergeSpannerDatabaseIamBindingDelete,
	}
}

func ResourceConverterSpannerDatabaseIamMember() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "spanner.googleapis.com/Database",
		Convert:           GetSpannerDatabaseIamMemberCaiObject,
		FetchFullResource: FetchSpannerDatabaseIamPolicy,
		MergeCreateUpdate: MergeSpannerDatabaseIamMember,
		MergeDelete:       MergeSpannerDatabaseIamMemberDelete,
	}
}

func GetSpannerDatabaseIamPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newSpannerDatabaseIamAsset(d, config, cai.ExpandIamPolicyBindings)
}

func GetSpannerDatabaseIamBindingCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newSpannerDatabaseIamAsset(d, config, cai.ExpandIamRoleBindings)
}

func GetSpannerDatabaseIamMemberCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newSpannerDatabaseIamAsset(d, config, cai.ExpandIamMemberBindings)
}

func MergeSpannerDatabaseIamPolicy(existing, incoming cai.Asset) cai.Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeSpannerDatabaseIamBinding(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeIamAssets(existing, incoming, cai.MergeAuthoritativeBindings)
}

func MergeSpannerDatabaseIamBindingDelete(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeDeleteIamAssets(existing, incoming, cai.MergeDeleteAuthoritativeBindings)
}

func MergeSpannerDatabaseIamMember(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeIamAssets(existing, incoming, cai.MergeAdditiveBindings)
}

func MergeSpannerDatabaseIamMemberDelete(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeDeleteIamAssets(existing, incoming, cai.MergeDeleteAdditiveBindings)
}

func newSpannerDatabaseIamAsset(
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	expandBindings func(d tpgresource.TerraformResourceData) ([]cai.IAMBinding, error),
) ([]cai.Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []cai.Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := cai.AssetName(d, config, "//spanner.googleapis.com/projects/{{project}}/instances/{{instance}}/databases/{{database}}")
	if err != nil {
		return []cai.Asset{}, err
	}

	return []cai.Asset{{
		Name: name,
		Type: "spanner.googleapis.com/Database",
		IAMPolicy: &cai.IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchSpannerDatabaseIamPolicy(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (cai.Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("instance"); !ok {
		return cai.Asset{}, cai.ErrEmptyIdentityField
	}

	if _, ok := d.GetOk("database"); !ok {
		return cai.Asset{}, cai.ErrEmptyIdentityField
	}

	return cai.FetchIamPolicy(
		NewSpannerDatabaseIamUpdater,
		d,
		config,
		"//spanner.googleapis.com/projects/{{project}}/instances/{{instance}}/databases/{{database}}",
		"spanner.googleapis.com/Database",
	)
}
