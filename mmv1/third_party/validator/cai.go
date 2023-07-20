package google

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
)

//Force build
type ConvertFunc = tpgresource.ConvertFunc
type GetApiObjectFunc = tpgresource.GetApiObjectFunc
type FetchFullResourceFunc = tpgresource.FetchFullResourceFunc
type MergeFunc = tpgresource.MergeFunc
type ResourceConverter = tpgresource.ResourceConverter

// Asset is the CAI representation of a resource.
type Asset = tpgresource.Asset

// AssetResource is the Asset's Resource field.
type AssetResource = tpgresource.AssetResource

// AssetName templates an asset.name by looking up and replacing all instances
// of {{field}}. In the case where a field would resolve to an empty string, a
// generated unique string will be used: "placeholder-" + randomString().
// This is done to preserve uniqueness of asset.name for a given asset.asset_type.
func AssetName(d tpgresource.TerraformResourceData, config *transport_tpg.Config, linkTmpl string) (string, error) {
	return tpgresource.AssetName(d, config, linkTmpl)
}

type Folder = tpgresource.Folder

type IAMPolicy = tpgresource.IAMPolicy

type IAMBinding = tpgresource.IAMBinding

type OrgPolicy = tpgresource.OrgPolicy

// V2OrgPolicies is the represtation of V2OrgPolicies
type V2OrgPolicies = tpgresource.V2OrgPolicies

// Spec is the representation of Spec for V2OrgPolicy
type PolicySpec = tpgresource.PolicySpec

type PolicyRule = tpgresource.PolicyRule

type StringValues = tpgresource.StringValues

type Expr = tpgresource.Expr

type Timestamp = tpgresource.Timestamp

type ListPolicyAllValues = tpgresource.ListPolicyAllValues

type ListPolicy = tpgresource.ListPolicy

type BooleanPolicy = tpgresource.BooleanPolicy
type RestoreDefault = tpgresource.RestoreDefault

func RandString(n int) string {
	return tpgresource.RandString(n)
}
