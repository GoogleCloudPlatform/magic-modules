package google

import (
	"fmt"
	"strings"
)

func resourceConverterKmsKeyRingIamPolicy() ResourceConverter {
	return ResourceConverter{
		AssetType:         "cloudkms.googleapis.com/KeyRing",
		Convert:           GetKmsKeyRingIamPolicyCaiObject,
		MergeCreateUpdate: MergeKmsKeyRingIamPolicy,
	}
}

func resourceConverterKmsKeyRingIamBinding() ResourceConverter {
	return ResourceConverter{
		AssetType:         "cloudkms.googleapis.com/KeyRing",
		Convert:           GetKmsKeyRingIamBindingCaiObject,
		FetchFullResource: FetchKmsKeyRingIamPolicy,
		MergeCreateUpdate: MergeKmsKeyRingIamBinding,
		MergeDelete:       MergeKmsKeyRingIamBindingDelete,
	}
}

func resourceConverterKmsKeyRingIamMember() ResourceConverter {
	return ResourceConverter{
		AssetType:         "cloudkms.googleapis.com/KeyRing",
		Convert:           GetKmsKeyRingIamMemberCaiObject,
		FetchFullResource: FetchKmsKeyRingIamPolicy,
		MergeCreateUpdate: MergeKmsKeyRingIamMember,
		MergeDelete:       MergeKmsKeyRingIamMemberDelete,
	}
}

func GetKmsKeyRingIamPolicyCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newKmsKeyRingIamAsset(d, config, expandIamPolicyBindings)
}

func GetKmsKeyRingIamBindingCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newKmsKeyRingIamAsset(d, config, expandIamRoleBindings)
}

func GetKmsKeyRingIamMemberCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	return newKmsKeyRingIamAsset(d, config, expandIamMemberBindings)
}

func MergeKmsKeyRingIamPolicy(existing, incoming Asset) Asset {
	existing.IAMPolicy = incoming.IAMPolicy
	return existing
}

func MergeKmsKeyRingIamBinding(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAuthoritativeBindings)
}

func MergeKmsKeyRingIamBindingDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAuthoritativeBindings)
}

func MergeKmsKeyRingIamMember(existing, incoming Asset) Asset {
	return mergeIamAssets(existing, incoming, mergeAdditiveBindings)
}

func MergeKmsKeyRingIamMemberDelete(existing, incoming Asset) Asset {
	return mergeDeleteIamAssets(existing, incoming, mergeDeleteAdditiveBindings)
}

func newKmsKeyRingIamAsset(
	d TerraformResourceData,
	config *Config,
	expandBindings func(d TerraformResourceData) ([]IAMBinding, error),
) ([]Asset, error) {
	bindings, err := expandBindings(d)
	if err != nil {
		return []Asset{}, fmt.Errorf("expanding bindings: %v", err)
	}

	assetNameTemplate := constructKmsKeyRingIAMAssetNameTemplate(d)
	name, err := assetName(d, config, assetNameTemplate)
	if err != nil {
		return []Asset{}, err
	}

	return []Asset{{
		Name: name,
		Type: "cloudkms.googleapis.com/KeyRing",
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}}, nil
}

func FetchKmsKeyRingIamPolicy(d TerraformResourceData, config *Config) (Asset, error) {
	// Check if the identity field returns a value
	if _, ok := d.GetOk("key_ring_id"); !ok {
		return Asset{}, ErrEmptyIdentityField
	}

	assetNameTemplate := constructKmsKeyRingIAMAssetNameTemplate(d)

	// We use key_ring_id in the asset name template to be consistent with newKmsKeyRingIamAsset.
	return fetchIamPolicy(
		NewKmsKeyRingIamUpdater,
		d,
		config,
		assetNameTemplate,                 // asset name
		"cloudkms.googleapis.com/KeyRing", // asset type
	)
}

func constructKmsKeyRingIAMAssetNameTemplate(d TerraformResourceData) string {
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
