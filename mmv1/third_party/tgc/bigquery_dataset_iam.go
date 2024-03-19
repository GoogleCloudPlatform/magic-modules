package bigquery

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

func ResourceConverterBigqueryDatasetIamPolicy() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "bigquery.googleapis.com/Dataset",
		Convert:           GetBigqueryDatasetIamPolicyCaiObject,
		MergeCreateUpdate: MergeBigqueryDatasetIamPolicy,
	}
}

func ResourceConverterBigqueryDatasetIamBinding() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "bigquery.googleapis.com/Dataset",
		Convert:           GetBigqueryDatasetIamBindingCaiObject,
		FetchFullResource: FetchBigqueryDatasetIamPolicy,
		MergeCreateUpdate: MergeBigqueryDatasetIamBinding,
		MergeDelete:       MergeBigqueryDatasetIamBindingDelete,
	}
}

func ResourceConverterBigqueryDatasetIamMember() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "bigquery.googleapis.com/Dataset",
		Convert:           GetBigqueryDatasetIamMemberCaiObject,
		FetchFullResource: FetchBigqueryDatasetIamPolicy,
		MergeCreateUpdate: MergeBigqueryDatasetIamMember,
		MergeDelete:       MergeBigqueryDatasetIamMemberDelete,
	}
}

func GetBigqueryDatasetIamPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newBigqueryDatasetIamAsset(d, config, cai.ExpandIamPolicyBindings)
}

func GetBigqueryDatasetIamBindingCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newBigqueryDatasetIamAsset(d, config, cai.ExpandIamRoleBindings)
}

func GetBigqueryDatasetIamMemberCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newBigqueryDatasetIamAsset(d, config, cai.ExpandIamMemberBindings)
}

func MergeBigqueryDatasetIamPolicy(existing, incoming cai.Asset) cai.Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeBigqueryDatasetIamBinding(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeIamAssets(existing, incoming, cai.MergeAuthoritativeBindings)
}

func MergeBigqueryDatasetIamBindingDelete(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeDeleteIamAssets(existing, incoming, cai.MergeDeleteAuthoritativeBindings)
}

func MergeBigqueryDatasetIamMember(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeIamAssets(existing, incoming, cai.MergeAdditiveBindings)
}

func MergeBigqueryDatasetIamMemberDelete(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeDeleteIamAssets(existing, incoming, cai.MergeDeleteAdditiveBindings)
}

func newBigqueryDatasetIamAsset(
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	expandBindings func(d tpgresource.TerraformResourceData) ([]cai.IAMBinding, error),
) ([]cai.Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []cai.Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	name, err := cai.AssetName(d, config, "//bigquery.googleapis.com/projects/{{project}}/datasets/{{dataset_id}}")
	if err != nil {
		return []cai.Asset{}, err
	}

	return []cai.Asset{{
		Name: name,
		Type: "bigquery.googleapis.com/Dataset",
		IAMPolicy: &cai.IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchBigqueryDatasetIamPolicy(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (cai.Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("dataset_id"); !ok {
		return cai.Asset{}, cai.ErrEmptyIdentityField
	}

	return cai.FetchIamPolicy(
		NewBigqueryDatasetIamUpdater,
		d,
		config,
		"//bigquery.googleapis.com/projects/{{project}}/datasets/{{dataset_id}}",
		"bigquery.googleapis.com/Dataset",
	)
}
