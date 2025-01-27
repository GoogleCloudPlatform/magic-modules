package google

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

type ConvertFunc = cai.ConvertFunc
type GetApiObjectFunc = cai.GetApiObjectFunc
type FetchFullResourceFunc = cai.FetchFullResourceFunc
type MergeFunc = cai.MergeFunc
type ResourceConverter = cai.ResourceConverter

// Asset is the CAI representation of a resource.
type Asset = cai.Asset

// AssetResource is the Asset's Resource field.
type AssetResource = cai.AssetResource

// AssetName templates an asset.name by looking up and replacing all instances
// of {{field}}. In the case where a field would resolve to an empty string, a
// generated unique string will be used: "placeholder-" + randomString().
// This is done to preserve uniqueness of asset.name for a given asset.asset_type.
func AssetName(d tpgresource.TerraformResourceData, config *transport_tpg.Config, linkTmpl string) (string, error) {
	return cai.AssetName(d, config, linkTmpl)
}

type Folder = cai.Folder

type IAMPolicy = cai.IAMPolicy

type IAMBinding = cai.IAMBinding

type OrgPolicy = cai.OrgPolicy

// V2OrgPolicies is the represtation of V2OrgPolicies
type V2OrgPolicies = cai.V2OrgPolicies

// Spec is the representation of Spec for V2OrgPolicy
type PolicySpec = cai.PolicySpec

type PolicyRule = cai.PolicyRule

type StringValues = cai.StringValues

type Expr = cai.Expr

type Timestamp = cai.Timestamp

type ListPolicyAllValues = cai.ListPolicyAllValues

type ListPolicy = cai.ListPolicy

type BooleanPolicy = cai.BooleanPolicy
type RestoreDefault = cai.RestoreDefault

func RandString(n int) string {
	return cai.RandString(n)
}
