// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tpgiamresource

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// ConvertToIdentitySchema maps an IAM parent resource schema to Terraform resource identity
// attributes (types and RequiredForImport / OptionalForImport). Used when building the identity
// schema for google_*_iam_member and related resources.
func ConvertToIdentitySchema(parentSchema map[string]*schema.Schema) map[string]*schema.Schema {
	identitySchema := make(map[string]*schema.Schema)
	for k, v := range parentSchema {
		identitySchema[k] = &schema.Schema{
			Type: v.Type,
		}
		// If the field has RequiredForImport or OptionalForImport set, preserve them
		if v.RequiredForImport {
			identitySchema[k].RequiredForImport = true
		}
		if v.OptionalForImport {
			identitySchema[k].OptionalForImport = true
		}
		// If not explicitly set, infer from Required+ForceNew pattern
		if !v.RequiredForImport && !v.OptionalForImport {
			if v.Required && v.ForceNew {
				identitySchema[k].RequiredForImport = true
			} else if v.Optional && v.ForceNew {
				identitySchema[k].OptionalForImport = true
			}
		}
	}
	return identitySchema
}

// PopulateIamParentIdentity copies IAM parent attributes from resource state into Terraform
// resource identity. Keys match ConvertToIdentitySchema(parentSchema); each
// ResourceIdentityParser used for import should read the same attributes when producing the
// canonical resource id (first segment of the legacy import id).
//
// This is shared by IAM fine-grained resources whose parent keys come from the merged
// Iam*Schema (project, folder, zone, name, etc.).
func PopulateIamParentIdentity(identity *schema.IdentityData, d *schema.ResourceData, parentSchema map[string]*schema.Schema) {
	for attr := range parentSchema {
		identity.Set(attr, d.Get(attr))
	}
}
