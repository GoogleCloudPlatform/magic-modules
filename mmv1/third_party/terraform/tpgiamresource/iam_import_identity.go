// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tpgiamresource

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// IamImportIdentityResourceAttributes reads the given top-level identity attribute names
// from Terraform import identity and returns them as a map of non-empty strings. The
// keys should match identity schema fields that correspond to named groups in the
// resource's import id regexes (e.g. project, zone, name for compute disks).
func IamImportIdentityResourceAttributes(identity *schema.IdentityData, keys []string) (map[string]string, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("identity attribute keys must not be empty")
	}
	out := make(map[string]string, len(keys))
	for _, k := range keys {
		v, ok := identity.GetOk(k)
		if !ok {
			return nil, fmt.Errorf("import identity is missing attribute %q", k)
		}
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("import identity attribute %q must be a string, got %T", k, v)
		}
		if s == "" {
			return nil, fmt.Errorf("import identity attribute %q is empty", k)
		}
		out[k] = s
	}
	return out, nil
}

// FormatIAMResourceCanonicalID builds the first segment of google_*_iam_member IDs — the
// same shape as ResourceIamUpdater.GetResourceId() — using a fmt.Sprintf-style format and
// attribute values in the order given by attributeKeys.
//
// Example (compute disk): format "projects/%s/zones/%s/disks/%s", keys []string{"project", "zone", "name"}.
func FormatIAMResourceCanonicalID(format string, attributeKeys []string, attrs map[string]string) (string, error) {
	if len(attributeKeys) == 0 {
		return "", fmt.Errorf("attributeKeys must not be empty")
	}
	if strings.Count(format, "%s") != len(attributeKeys) {
		return "", fmt.Errorf("format %q must contain exactly %d %%s placeholders (one per attribute key)", format, len(attributeKeys))
	}
	args := make([]interface{}, 0, len(attributeKeys))
	for _, k := range attributeKeys {
		v, ok := attrs[k]
		if !ok {
			return "", fmt.Errorf("attrs map is missing key %q", k)
		}
		args = append(args, v)
	}
	return fmt.Sprintf(format, args...), nil
}

// CanonicalResourceIDFromIamImportIdentity combines IamImportIdentityResourceAttributes and
// FormatIAMResourceCanonicalID for resource-identity import of IAM member resources.
func CanonicalResourceIDFromIamImportIdentity(identity *schema.IdentityData, attributeKeys []string, canonicalFormat string) (string, error) {
	attrs, err := IamImportIdentityResourceAttributes(identity, attributeKeys)
	if err != nil {
		return "", err
	}
	return FormatIAMResourceCanonicalID(canonicalFormat, attributeKeys, attrs)
}
