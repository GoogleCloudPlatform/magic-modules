package cai

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/caiasset"

	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

type ConvertFunc func(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]caiasset.Asset, error)

// FetchFullResourceFunc allows initial data for a resource to be fetched from the API and merged
// with the planned changes. This is useful for resources that are only partially managed
// by Terraform, like IAM policies managed with member/binding resources.
type FetchFullResourceFunc func(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (caiasset.Asset, error)

type ResourceConverter struct {
	Convert           ConvertFunc
	FetchFullResource FetchFullResourceFunc
}
