package google

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/services/kms"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgiamresource"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
)

func resourceConverterKmsCryptoKeyIamPolicy() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "cloudkms.googleapis.com/CryptoKey",
		Convert:           GetKmsCryptoKeyIamPolicyCaiObject,
		MergeCreateUpdate: MergeKmsCryptoKeyIamPolicy,
	}
}

func resourceConverterKmsCryptoKeyIamBinding() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "cloudkms.googleapis.com/CryptoKey",
		Convert:           GetKmsCryptoKeyIamBindingCaiObject,
		FetchFullResource: FetchKmsCryptoKeyIamPolicy,
		MergeCreateUpdate: MergeKmsCryptoKeyIamBinding,
		MergeDelete:       MergeKmsCryptoKeyIamBindingDelete,
	}
}

func resourceConverterKmsCryptoKeyIamMember() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "cloudkms.googleapis.com/CryptoKey",
		Convert:           GetKmsCryptoKeyIamMemberCaiObject,
		FetchFullResource: FetchKmsCryptoKeyIamPolicy,
		MergeCreateUpdate: MergeKmsCryptoKeyIamMember,
		MergeDelete:       MergeKmsCryptoKeyIamMemberDelete,
	}
}

func GetKmsCryptoKeyIamPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newKmsCryptoKeyIamAsset(d, config, tpgiamresource.ExpandIamPolicyBindings)
}

func GetKmsCryptoKeyIamBindingCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newKmsCryptoKeyIamAsset(d, config, tpgiamresource.ExpandIamRoleBindings)
}

func GetKmsCryptoKeyIamMemberCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	return newKmsCryptoKeyIamAsset(d, config, tpgiamresource.ExpandIamMemberBindings)
}

func MergeKmsCryptoKeyIamPolicy(existing, incoming tpgresource.Asset) tpgresource.Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeKmsCryptoKeyIamBinding(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeIamAssets(existing, incoming, tpgiamresource.MergeAuthoritativeBindings)
}

func MergeKmsCryptoKeyIamBindingDelete(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeDeleteIamAssets(existing, incoming, tpgiamresource.MergeDeleteAuthoritativeBindings)
}

func MergeKmsCryptoKeyIamMember(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeIamAssets(existing, incoming, tpgiamresource.MergeAdditiveBindings)
}

func MergeKmsCryptoKeyIamMemberDelete(existing, incoming tpgresource.Asset) tpgresource.Asset {
	return tpgiamresource.MergeDeleteIamAssets(existing, incoming, tpgiamresource.MergeDeleteAdditiveBindings)
}

func newKmsCryptoKeyIamAsset(
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	expandBindings func(d tpgresource.TerraformResourceData) ([]tpgresource.IAMBinding, error),
) ([]tpgresource.Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []tpgresource.Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	assetNameTemplate := constructAssetNameTemplate(d)
	name, err := tpgresource.AssetName(d, config, assetNameTemplate)
	if err != nil {
		return []tpgresource.Asset{}, err
	}

	return []tpgresource.Asset{{
		Name: name,
		Type: "cloudkms.googleapis.com/CryptoKey",
		IAMPolicy: &tpgresource.IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchKmsCryptoKeyIamPolicy(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgresource.Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("crypto_key_id"); !ok {
		return tpgresource.Asset{}, tpgresource.ErrEmptyIdentityField
	}

	assetNameTemplate := constructAssetNameTemplate(d)

	// We use crypto_key_id in the asset name template to be consistent with newKmsCryptoKeyIamAsset.
	return tpgiamresource.FetchIamPolicy(
		kms.NewKmsCryptoKeyIamUpdater,
		d,
		config,
		assetNameTemplate,                   // asset name
		"cloudkms.googleapis.com/CryptoKey", // asset type
	)
}

func constructAssetNameTemplate(d tpgresource.TerraformResourceData) string {
	assetNameTemplate := "//cloudkms.googleapis.com/{{crypto_key_id}}"
	if val, ok := d.GetOk("crypto_key_id"); ok {
		cryptoKeyID := val.(string)
		splits := strings.Split(cryptoKeyID, "/")
		if len(splits) == 4 {
			assetNameTemplate = fmt.Sprintf("//cloudkms.googleapis.com/projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", splits[0], splits[1], splits[2], splits[3])
		} else if len(splits) == 3 {
			assetNameTemplate = fmt.Sprintf("//cloudkms.googleapis.com/projects/{{project}}/locations/%s/keyRings/%s/cryptoKeys/%s", splits[0], splits[1], splits[2])
		}
	}
	return assetNameTemplate
}
