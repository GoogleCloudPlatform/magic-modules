package compute

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

// updateKmsKey can't remove a key or use a key version, clears
// kmsKeyServiceAccount, and rotates if kmsKeyName is empty.

// kmsKeyChange reports whether a key change needs replacement, or an error if
// it can't be made. serviceAccountSet is true if kms_key_service_account is
// set; customerSuppliedKey is true if the resource uses or switches to a CSEK.
func kmsKeyChange(newKey string, newKeyKnown, serviceAccountSet, customerSuppliedKey bool) (bool, error) {
	if serviceAccountSet || customerSuppliedKey {
		return true, nil
	}
	if !newKeyKnown {
		// Once known, this is an update or an error, never a replacement.
		return false, nil
	}
	if newKey == "" {
		// -replace still hits this error; taint doesn't.
		return false, fmt.Errorf("removing the Cloud KMS key isn't supported; to recreate the resource without a key, taint or destroy it first")
	}
	if isCryptoKeyVersionName(newKey) {
		return false, fmt.Errorf("Cloud KMS key %q includes a crypto key version, which isn't supported; remove the /cryptoKeyVersions/ suffix", newKey)
	}
	return false, nil
}

func isCryptoKeyVersionName(key string) bool {
	return strings.Contains(key, "/cryptoKeyVersions/")
}

// validateKmsKeyChange forces replacement or errors at plan time for key
// changes updateKmsKey can't make. serviceAccountPath may be "".
func validateKmsKeyChange(keyPath, serviceAccountPath, sha256Path string) schema.CustomizeDiffFunc {
	return func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
		if d.Id() == "" || !d.HasChange(keyPath) {
			return nil
		}
		oldKey, newKey := d.GetChange(keyPath)
		// HasChange ignores DiffSuppressFunc.
		if d.NewValueKnown(keyPath) && tpgresource.CompareKmsKeyNames(keyPath, oldKey.(string), newKey.(string), nil) {
			return nil
		}
		serviceAccountSet := false
		if serviceAccountPath != "" {
			oldServiceAccount, newServiceAccount := d.GetChange(serviceAccountPath)
			serviceAccountSet = oldServiceAccount.(string) != "" || newServiceAccount.(string) != ""
		}
		oldSha256, _ := d.GetChange(sha256Path)
		customerSuppliedKey := (oldKey.(string) == "" && oldSha256.(string) != "") ||
			customerSuppliedKeyInConfig(d, strings.SplitN(keyPath, ".", 2)[0])
		forceNew, err := kmsKeyChange(newKey.(string), d.NewValueKnown(keyPath), serviceAccountSet, customerSuppliedKey)
		if err != nil {
			return err
		}
		if forceNew {
			return d.ForceNew(keyPath)
		}
		return nil
	}
}

// customerSuppliedKeyInConfig reports whether the config sets a raw or
// RSA-wrapped key, including write-only variants, in the encryption block.
func customerSuppliedKeyInConfig(d *schema.ResourceDiff, block string) bool {
	for _, f := range []string{"raw_key", "rsa_encrypted_key"} {
		path := block + ".0." + f
		if v, ok := d.Get(path).(string); ok && v != "" || !d.NewValueKnown(path) {
			return true
		}
	}
	rawConfig := d.GetRawConfig()
	if !rawConfig.IsKnown() || rawConfig.IsNull() || !rawConfig.Type().IsObjectType() || !rawConfig.Type().HasAttribute(block) {
		return false
	}
	blockValue := rawConfig.GetAttr(block)
	if !blockValue.IsKnown() || blockValue.IsNull() || !blockValue.CanIterateElements() || blockValue.LengthInt() == 0 {
		return false
	}
	for it := blockValue.ElementIterator(); it.Next(); {
		_, elem := it.Element()
		for _, f := range []string{"raw_key_wo", "rsa_encrypted_key_wo"} {
			if !elem.Type().IsObjectType() || !elem.Type().HasAttribute(f) {
				continue
			}
			if v := elem.GetAttr(f); !v.IsNull() && (!v.IsKnown() || v.AsString() != "") {
				return true
			}
		}
	}
	return false
}

func kmsKeyUpdateRequestBody(newKey, serviceAccount string) (map[string]interface{}, error) {
	switch {
	case newKey == "":
		return nil, fmt.Errorf("removing the Cloud KMS key isn't supported")
	case isCryptoKeyVersionName(newKey):
		return nil, fmt.Errorf("the Cloud KMS key must not include a crypto key version (%q)", newKey)
	case serviceAccount != "":
		return nil, fmt.Errorf("the Cloud KMS key can't be changed in place while kms_key_service_account is set")
	}
	return map[string]interface{}{"kmsKeyName": newKey}, nil
}
