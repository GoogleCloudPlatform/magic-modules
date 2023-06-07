package google

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgiamresource"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
)

func resourceConverterSpannerDatabaseIamPolicy() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "spanner.googleapis.com/Database",
		Convert:           GetSpannerDatabaseIamPolicyCaiObject,
		MergeCreateUpdate: MergeSpannerDatabaseIamPolicy,
	}
}

func resourceConverterSpannerDatabaseIamBinding() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "spanner.googleapis.com/Database",
		Convert:           GetSpannerDatabaseIamBindingCaiObject,
		FetchFullResource: FetchSpannerDatabaseIamPolicy,
		MergeCreateUpdate: MergeSpannerDatabaseIamBinding,
		MergeDelete:       MergeSpannerDatabaseIamBindingDelete,
	}
}

func resourceConverterSpannerDatabaseIamMember() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "spanner.googleapis.com/Database",
		Convert:           GetSpannerDatabaseIamMemberCaiObject,
		FetchFullResource: FetchSpannerDatabaseIamPolicy,
		MergeCreateUpdate: MergeSpannerDatabaseIamMember,
		MergeDelete:       MergeSpannerDatabaseIamMemberDelete,
	}
}

func GetSpannerDatabaseIamPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newSpannerDatabaseIamAsset(d, config, tpgiamresource.ExpandIamPolicyBindings)
}

func GetSpannerDatabaseIamBindingCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newSpannerDatabaseIamAsset(d, config, tpgiamresource.ExpandIamRoleBindings)
}

func GetSpannerDatabaseIamMemberCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newSpannerDatabaseIamAsset(d, config, tpgiamresource.ExpandIamMemberBindings)
}

func MergeSpannerDatabaseIamPolicy(existing, incoming tpgresource.Asset) tpgresource.Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeSpannerDatabaseIamBinding(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeIamAssets(existing, incoming, tpgiamresource.MergeAuthoritativeBindings)
}

func MergeSpannerDatabaseIamBindingDelete(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeDeleteIamAssets(existing, incoming, tpgiamresource.MergeDeleteAuthoritativeBindings)
}

func MergeSpannerDatabaseIamMember(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeIamAssets(existing, incoming, tpgiamresource.MergeAdditiveBindings)
}

func MergeSpannerDatabaseIamMemberDelete(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeDeleteIamAssets(existing, incoming, tpgiamresource.MergeDeleteAdditiveBindings)
}

func newSpannerDatabaseIamAsset(
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	expandBindings func(d tpgresource.TerraformResourceData) ([]tpgresource.IAMBinding, error),
) ([]tpgresource.Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []tpgresource.Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := tpgresource.AssetName(d, config, "//spanner.googleapis.com/projects/{{project}}/instances/{{instance}}/databases/{{database}}")
	if err != nil {
		return []tpgresource.Asset{}, err
	}

	return []tpgresource.Asset{{
		Name: name,
		Type: "spanner.googleapis.com/Database",
		IAMPolicy: &tpgresource.IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchSpannerDatabaseIamPolicy(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgresource.Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("instance"); !ok {
		return tpgresource.Asset{}, tpgresource.ErrEmptyIdentityField
	}

	if _, ok := d.GetOk("database"); !ok {
		return tpgresource.Asset{}, tpgresource.ErrEmptyIdentityField
	}

	return tpgiamresource.FetchIamPolicy(
		NewSpannerDatabaseIamUpdater,
		d,
		config,
		"//spanner.googleapis.com/projects/{{project}}/instances/{{instance}}/databases/{{database}}",
		"spanner.googleapis.com/Database",
	)
}
