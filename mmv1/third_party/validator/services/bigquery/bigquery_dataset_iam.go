package bigquery

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgiamresource"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
)

func ResourceConverterBigqueryDatasetIamPolicy() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "bigquery.googleapis.com/Dataset",
		Convert:           GetBigqueryDatasetIamPolicyCaiObject,
		MergeCreateUpdate: MergeBigqueryDatasetIamPolicy,
	}
}

func ResourceConverterBigqueryDatasetIamBinding() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "bigquery.googleapis.com/Dataset",
		Convert:           GetBigqueryDatasetIamBindingCaiObject,
		FetchFullResource: FetchBigqueryDatasetIamPolicy,
		MergeCreateUpdate: MergeBigqueryDatasetIamBinding,
		MergeDelete:       MergeBigqueryDatasetIamBindingDelete,
	}
}

func ResourceConverterBigqueryDatasetIamMember() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "bigquery.googleapis.com/Dataset",
		Convert:           GetBigqueryDatasetIamMemberCaiObject,
		FetchFullResource: FetchBigqueryDatasetIamPolicy,
		MergeCreateUpdate: MergeBigqueryDatasetIamMember,
		MergeDelete:       MergeBigqueryDatasetIamMemberDelete,
	}
}

func GetBigqueryDatasetIamPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newBigqueryDatasetIamAsset(d, config, tpgiamresource.ExpandIamPolicyBindings)
}

func GetBigqueryDatasetIamBindingCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newBigqueryDatasetIamAsset(d, config, tpgiamresource.ExpandIamRoleBindings)
}

func GetBigqueryDatasetIamMemberCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newBigqueryDatasetIamAsset(d, config, tpgiamresource.ExpandIamMemberBindings)
}

func MergeBigqueryDatasetIamPolicy(existing, incoming tpgresource.Asset) tpgresource.Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeBigqueryDatasetIamBinding(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeIamAssets(existing, incoming, tpgiamresource.MergeAuthoritativeBindings)
}

func MergeBigqueryDatasetIamBindingDelete(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeDeleteIamAssets(existing, incoming, tpgiamresource.MergeDeleteAuthoritativeBindings)
}

func MergeBigqueryDatasetIamMember(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeIamAssets(existing, incoming, tpgiamresource.MergeAdditiveBindings)
}

func MergeBigqueryDatasetIamMemberDelete(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeDeleteIamAssets(existing, incoming, tpgiamresource.MergeDeleteAdditiveBindings)
}

func newBigqueryDatasetIamAsset(
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	expandBindings func(d tpgresource.TerraformResourceData) ([]tpgresource.IAMBinding, error),
) ([]tpgresource.Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []tpgresource.Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := tpgresource.AssetName(d, config, "//bigquery.googleapis.com/projects/{{project}}/datasets/{{dataset_id}}")
	if err != nil {
		return []tpgresource.Asset{}, err
	}

	return []tpgresource.Asset{{
		Name: name,
		Type: "bigquery.googleapis.com/Dataset",
		IAMPolicy: &tpgresource.IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchBigqueryDatasetIamPolicy(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgresource.Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("dataset_id"); !ok {
		return tpgresource.Asset{}, tpgresource.ErrEmptyIdentityField
	}

	return tpgiamresource.FetchIamPolicy(
		NewBigqueryDatasetIamUpdater,
		d,
		config,
		"//bigquery.googleapis.com/projects/{{project}}/datasets/{{dataset_id}}",
		"bigquery.googleapis.com/Dataset",
	)
}
