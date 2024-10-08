package kms

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/services/kms"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

func ResourceConverterKmsKeyRingIamPolicy() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "cloudkms.googleapis.com/KeyRing",
		Convert:           GetKmsKeyRingIamPolicyCaiObject,
		MergeCreateUpdate: MergeKmsKeyRingIamPolicy,
	}
}

func ResourceConverterKmsKeyRingIamBinding() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "cloudkms.googleapis.com/KeyRing",
		Convert:           GetKmsKeyRingIamBindingCaiObject,
		FetchFullResource: FetchKmsKeyRingIamPolicy,
		MergeCreateUpdate: MergeKmsKeyRingIamBinding,
		MergeDelete:       MergeKmsKeyRingIamBindingDelete,
	}
}

func ResourceConverterKmsKeyRingIamMember() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "cloudkms.googleapis.com/KeyRing",
		Convert:           GetKmsKeyRingIamMemberCaiObject,
		FetchFullResource: FetchKmsKeyRingIamPolicy,
		MergeCreateUpdate: MergeKmsKeyRingIamMember,
		MergeDelete:       MergeKmsKeyRingIamMemberDelete,
	}
}

func GetKmsKeyRingIamPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newKmsKeyRingIamAsset(d, config, cai.ExpandIamPolicyBindings)
}

func GetKmsKeyRingIamBindingCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newKmsKeyRingIamAsset(d, config, cai.ExpandIamRoleBindings)
}

func GetKmsKeyRingIamMemberCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newKmsKeyRingIamAsset(d, config, cai.ExpandIamMemberBindings)
}

func MergeKmsKeyRingIamPolicy(existing, incoming cai.Asset) cai.Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeKmsKeyRingIamBinding(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeIamAssets(existing, incoming, cai.MergeAuthoritativeBindings)
}

func MergeKmsKeyRingIamBindingDelete(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeDeleteIamAssets(existing, incoming, cai.MergeDeleteAuthoritativeBindings)
}

func MergeKmsKeyRingIamMember(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeIamAssets(existing, incoming, cai.MergeAdditiveBindings)
}

func MergeKmsKeyRingIamMemberDelete(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeDeleteIamAssets(existing, incoming, cai.MergeDeleteAdditiveBindings)
}

func newKmsKeyRingIamAsset(
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	expandBindings func(d tpgresource.TerraformResourceData) ([]cai.IAMBinding, error),
) ([]cai.Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []cai.Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	assetNameTemplate := constructKmsKeyRingIAMAssetNameTemplate(d)
	name, err := cai.AssetName(d, config, assetNameTemplate)
	if err != nil {
		return []cai.Asset{}, err
	}

	return []cai.Asset{{
		Name: name,
		Type: "cloudkms.googleapis.com/KeyRing",
		IAMPolicy: &cai.IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchKmsKeyRingIamPolicy(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (cai.Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("key_ring_id"); !ok {
		return cai.Asset{}, cai.ErrEmptyIdentityField
	}

	assetNameTemplate := constructKmsKeyRingIAMAssetNameTemplate(d)

	// We use key_ring_id in the asset name template to be consistent with newKmsKeyRingIamAsset.
	return cai.FetchIamPolicy(
		kms.NewKmsKeyRingIamUpdater,
		d,
		config,
		assetNameTemplate,                 // asset name
		"cloudkms.googleapis.com/KeyRing", // asset type
	)
}

func constructKmsKeyRingIAMAssetNameTemplate(d tpgresource.TerraformResourceData) string {
	assetNameTemplate := "//cloudkms.googleapis.com/{{key_ring_id}}"
	if val, ok := d.GetOk("key_ring_id"); ok {
		keyRingID := val.(string)
		splits := strings.Split(keyRingID, "/")
		if len(splits) == 3 {
			assetNameTemplate = fmt.Sprintf("//cloudkms.googleapis.com/projects/%s/locations/%s/keyRings/%s", splits[0], splits[1], splits[2])
		} else if len(splits) == 2 {
			assetNameTemplate = fmt.Sprintf("//cloudkms.googleapis.com/projects/{{project}}/locations/%s/keyRings/%s", splits[0], splits[1])
		}
	}
	return assetNameTemplate
}
