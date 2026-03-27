// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tpgresource

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// IamImportIdentityData wraps import-time identity so GetProject, GetRegion, GetZone, and
// GetLocation can resolve provider defaults when optional parent attributes are omitted,
// consistent with IAM updater production and tpgresource.GetImportIdQualifiers defaulting.
type IamImportIdentityData struct {
	Identity *schema.IdentityData
}

// IamImportTerraformData returns identity as TerraformResourceData for location/project helpers.
// identity must be non-nil when passed to GetProject / GetRegion / GetZone / GetLocation.
func IamImportTerraformData(identity *schema.IdentityData) TerraformResourceData {
	return &IamImportIdentityData{Identity: identity}
}

func (w *IamImportIdentityData) HasChange(string) bool { return false }

func (w *IamImportIdentityData) GetOkExists(key string) (interface{}, bool) {
	return w.GetOk(key)
}

func (w *IamImportIdentityData) GetOk(key string) (interface{}, bool) {
	if w.Identity == nil {
		return nil, false
	}
	return w.Identity.GetOk(key)
}

func (w *IamImportIdentityData) Get(key string) interface{} {
	if w.Identity == nil {
		return nil
	}
	return w.Identity.Get(key)
}

func (w *IamImportIdentityData) Set(key string, value interface{}) error {
	if w.Identity == nil {
		return nil
	}
	return w.Identity.Set(key, value)
}

func (w *IamImportIdentityData) SetId(string) {}

func (w *IamImportIdentityData) Id() string { return "" }

func (w *IamImportIdentityData) GetProviderMeta(interface{}) error { return nil }

func (w *IamImportIdentityData) Timeout(string) time.Duration { return 0 }
