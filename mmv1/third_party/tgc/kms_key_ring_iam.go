package google

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/services/kms"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgiamresource"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
)

func resourceConverterKmsKeyRingIamPolicy() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "cloudkms.googleapis.com/KeyRing",
		Convert:           GetKmsKeyRingIamPolicyCaiObject,
		MergeCreateUpdate: MergeKmsKeyRingIamPolicy,
	}
}

func resourceConverterKmsKeyRingIamBinding() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "cloudkms.googleapis.com/KeyRing",
		Convert:           GetKmsKeyRingIamBindingCaiObject,
		FetchFullResource: FetchKmsKeyRingIamPolicy,
		MergeCreateUpdate: MergeKmsKeyRingIamBinding,
		MergeDelete:       MergeKmsKeyRingIamBindingDelete,
	}
}

func resourceConverterKmsKeyRingIamMember() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "cloudkms.googleapis.com/KeyRing",
		Convert:           GetKmsKeyRingIamMemberCaiObject,
		FetchFullResource: FetchKmsKeyRingIamPolicy,
		MergeCreateUpdate: MergeKmsKeyRingIamMember,
		MergeDelete:       MergeKmsKeyRingIamMemberDelete,
	}
}

func GetKmsKeyRingIamPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newKmsKeyRingIamAsset(d, config, tpgiamresource.ExpandIamPolicyBindings)
}

func GetKmsKeyRingIamBindingCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newKmsKeyRingIamAsset(d, config, tpgiamresource.ExpandIamRoleBindings)
}

func GetKmsKeyRingIamMemberCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newKmsKeyRingIamAsset(d, config, tpgiamresource.ExpandIamMemberBindings)
}

func MergeKmsKeyRingIamPolicy(existing, incoming tpgresource.Asset) tpgresource.Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeKmsKeyRingIamBinding(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeIamAssets(existing, incoming, tpgiamresource.MergeAuthoritativeBindings)
}

func MergeKmsKeyRingIamBindingDelete(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeDeleteIamAssets(existing, incoming, tpgiamresource.MergeDeleteAuthoritativeBindings)
}

func MergeKmsKeyRingIamMember(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeIamAssets(existing, incoming, tpgiamresource.MergeAdditiveBindings)
}

func MergeKmsKeyRingIamMemberDelete(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeDeleteIamAssets(existing, incoming, tpgiamresource.MergeDeleteAdditiveBindings)
}

func newKmsKeyRingIamAsset(
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	expandBindings func(d tpgresource.TerraformResourceData) ([]tpgresource.IAMBinding, error),
) ([]tpgresource.Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []tpgresource.Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	assetNameTemplate := constructKmsKeyRingIAMAssetNameTemplate(d)
	name, err := tpgresource.AssetName(d, config, assetNameTemplate)
	if err != nil {
		return []tpgresource.Asset{}, err
	}

	return []tpgresource.Asset{{
		Name: name,
		Type: "cloudkms.googleapis.com/KeyRing",
		IAMPolicy: &tpgresource.IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchKmsKeyRingIamPolicy(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgresource.Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("key_ring_id"); !ok {
		return tpgresource.Asset{}, tpgresource.ErrEmptyIdentityField
	}

	assetNameTemplate := constructKmsKeyRingIAMAssetNameTemplate(d)

	// We use key_ring_id in the asset name template to be consistent with newKmsKeyRingIamAsset.
	return tpgiamresource.FetchIamPolicy(
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
