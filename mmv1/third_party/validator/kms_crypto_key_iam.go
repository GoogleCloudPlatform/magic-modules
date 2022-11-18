package google

import (
	"fmt"
	"strings"
)

func resourceConverterKmsCryptoKeyIamPolicy() ResourceConverter {
	return ResourceConverter{
		AssetType:         "cloudkms.googleapis.com/CryptoKey",
		Convert:           GetKmsCryptoKeyIamPolicyCaiObject,
		MergeCreateUpdate: MergeKmsCryptoKeyIamPolicy,
	}
}

func resourceConverterKmsCryptoKeyIamBinding() ResourceConverter {
	return ResourceConverter{
		AssetType:         "cloudkms.googleapis.com/CryptoKey",
		Convert:           GetKmsCryptoKeyIamBindingCaiObject,
		FetchFullResource: FetchKmsCryptoKeyIamPolicy,
		MergeCreateUpdate: MergeKmsCryptoKeyIamBinding,
		MergeDelete:       MergeKmsCryptoKeyIamBindingDelete,
	}
}

func resourceConverterKmsCryptoKeyIamMember() ResourceConverter {
	return ResourceConverter{
		AssetType:         "cloudkms.googleapis.com/CryptoKey",
		Convert:           GetKmsCryptoKeyIamMemberCaiObject,
		FetchFullResource: FetchKmsCryptoKeyIamPolicy,
		MergeCreateUpdate: MergeKmsCryptoKeyIamMember,
		MergeDelete:       MergeKmsCryptoKeyIamMemberDelete,
	}
}

func GetKmsCryptoKeyIamPolicyCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newKmsCryptoKeyIamAsset(d, config, expandIamPolicyBindings)
}

func GetKmsCryptoKeyIamBindingCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newKmsCryptoKeyIamAsset(d, config, expandIamRoleBindings)
}

func GetKmsCryptoKeyIamMemberCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newKmsCryptoKeyIamAsset(d, config, expandIamMemberBindings)
}

func MergeKmsCryptoKeyIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeKmsCryptoKeyIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeKmsCryptoKeyIamBindingDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAuthoritativeBindings)
}

func MergeKmsCryptoKeyIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func MergeKmsCryptoKeyIamMemberDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAdditiveBindings)
}

func newKmsCryptoKeyIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) ([]Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	assetNameTemplate := constructAssetNameTemplate(d)
	name, err := assetName(d, config, assetNameTemplate)
	if err != nil {
		return []Asset{}, err
	}

	return []Asset{{
		Name: name,
		Type: "cloudkms.googleapis.com/CryptoKey",
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchKmsCryptoKeyIamPolicy(d TerraformResourceData, config *Config) (Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("crypto_key_id"); !ok {
		return Asset{}, ErrEmptyIdentityField
	}

	assetNameTemplate := constructAssetNameTemplate(d)

	// We use crypto_key_id in the asset name template to be consistent with newKmsCryptoKeyIamAsset.
	return fetchIamPolicy(
		NewKmsCryptoKeyIamUpdater,
		d,
		config,
		assetNameTemplate,                   // asset name
		"cloudkms.googleapis.com/CryptoKey", // asset type
	)
}

func constructAssetNameTemplate(d TerraformResourceData) string {
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
