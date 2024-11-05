package kms

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/services/kms"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

func ResourceConverterKmsCryptoKeyIamPolicy() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "cloudkms.googleapis.com/CryptoKey",
		Convert:           GetKmsCryptoKeyIamPolicyCaiObject,
		MergeCreateUpdate: MergeKmsCryptoKeyIamPolicy,
	}
}

func ResourceConverterKmsCryptoKeyIamBinding() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "cloudkms.googleapis.com/CryptoKey",
		Convert:           GetKmsCryptoKeyIamBindingCaiObject,
		FetchFullResource: FetchKmsCryptoKeyIamPolicy,
		MergeCreateUpdate: MergeKmsCryptoKeyIamBinding,
		MergeDelete:       MergeKmsCryptoKeyIamBindingDelete,
	}
}

func ResourceConverterKmsCryptoKeyIamMember() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType:         "cloudkms.googleapis.com/CryptoKey",
		Convert:           GetKmsCryptoKeyIamMemberCaiObject,
		FetchFullResource: FetchKmsCryptoKeyIamPolicy,
		MergeCreateUpdate: MergeKmsCryptoKeyIamMember,
		MergeDelete:       MergeKmsCryptoKeyIamMemberDelete,
	}
}

func GetKmsCryptoKeyIamPolicyCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newKmsCryptoKeyIamAsset(d, config, cai.ExpandIamPolicyBindings)
}

func GetKmsCryptoKeyIamBindingCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newKmsCryptoKeyIamAsset(d, config, cai.ExpandIamRoleBindings)
}

func GetKmsCryptoKeyIamMemberCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	return newKmsCryptoKeyIamAsset(d, config, cai.ExpandIamMemberBindings)
}

func MergeKmsCryptoKeyIamPolicy(existing, incoming cai.Asset) cai.Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeKmsCryptoKeyIamBinding(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeIamAssets(existing, incoming, cai.MergeAuthoritativeBindings)
}

func MergeKmsCryptoKeyIamBindingDelete(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeDeleteIamAssets(existing, incoming, cai.MergeDeleteAuthoritativeBindings)
}

func MergeKmsCryptoKeyIamMember(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeIamAssets(existing, incoming, cai.MergeAdditiveBindings)
}

func MergeKmsCryptoKeyIamMemberDelete(existing, incoming cai.Asset) cai.Asset {
	return cai.MergeDeleteIamAssets(existing, incoming, cai.MergeDeleteAdditiveBindings)
}

func newKmsCryptoKeyIamAsset(
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	expandBindings func(d tpgresource.TerraformResourceData) ([]cai.IAMBinding, error),
) ([]cai.Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []cai.Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	assetNameTemplate := constructAssetNameTemplate(d)
	name, err := cai.AssetName(d, config, assetNameTemplate)
	if err != nil {
		return []cai.Asset{}, err
	}

	return []cai.Asset{{
		Name: name,
		Type: "cloudkms.googleapis.com/CryptoKey",
		IAMPolicy: &cai.IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchKmsCryptoKeyIamPolicy(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (cai.Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("crypto_key_id"); !ok {
		return cai.Asset{}, cai.ErrEmptyIdentityField
	}

	assetNameTemplate := constructAssetNameTemplate(d)

	// We use crypto_key_id in the asset name template to be consistent with newKmsCryptoKeyIamAsset.
	return cai.FetchIamPolicy(
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
