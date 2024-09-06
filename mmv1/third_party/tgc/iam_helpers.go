package google

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

// ExpandIamPolicyBindings is used in google_<type>_iam_policy resources.
func expandIamPolicyBindings(d tpgresource.TerraformResourceData) ([]cai.IAMBinding, error) {
	return cai.ExpandIamPolicyBindings(d)
}

// ExpandIamRoleBindings is used in google_<type>_iam_binding resources.
func expandIamRoleBindings(d tpgresource.TerraformResourceData) ([]cai.IAMBinding, error) {
	return cai.ExpandIamRoleBindings(d)
}

// ExpandIamMemberBindings is used in google_<type>_iam_member resources.
func expandIamMemberBindings(d tpgresource.TerraformResourceData) ([]cai.IAMBinding, error) {
	return cai.ExpandIamMemberBindings(d)
}

// MergeIamAssets merges an existing asset with the IAM bindings of an incoming
// cai.Asset.
func mergeIamAssets(
	existing, incoming cai.Asset,
	MergeBindings func(existing, incoming []cai.IAMBinding) []cai.IAMBinding,
) cai.Asset {
	return cai.MergeIamAssets(
		existing, incoming,
		MergeBindings,
	)
}

// incoming is the last known state of an asset prior to deletion
func mergeDeleteIamAssets(
	existing, incoming cai.Asset,
	MergeBindings func(existing, incoming []cai.IAMBinding) []cai.IAMBinding,
) cai.Asset {
	return cai.MergeDeleteIamAssets(
		existing, incoming,
		MergeBindings,
	)
}

// MergeAdditiveBindings adds members to bindings with the same roles and adds new
// bindings for roles that dont exist.
func mergeAdditiveBindings(existing, incoming []cai.IAMBinding) []cai.IAMBinding {
	return cai.MergeAdditiveBindings(existing, incoming)
}

// MergeDeleteAdditiveBindings eliminates listed members from roles in the
// existing list. incoming is the last known state of the bindings being deleted.
func mergeDeleteAdditiveBindings(existing, incoming []cai.IAMBinding) []cai.IAMBinding {
	return cai.MergeDeleteAdditiveBindings(existing, incoming)
}

// MergeAuthoritativeBindings clobbers members to bindings with the same roles
// and adds new bindings for roles that dont exist.
func mergeAuthoritativeBindings(existing, incoming []cai.IAMBinding) []cai.IAMBinding {
	return cai.MergeAuthoritativeBindings(existing, incoming)
}

// MergeDeleteAuthoritativeBindings eliminates any bindings with matching roles
// in the existing list. incoming is the last known state of the bindings being
// deleted.
func mergeDeleteAuthoritativeBindings(existing, incoming []cai.IAMBinding) []cai.IAMBinding {
	return cai.MergeDeleteAuthoritativeBindings(existing, incoming)
}

func fetchIamPolicy(
	newUpdaterFunc tpgiamresource.NewResourceIamUpdaterFunc,
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	assetNameTmpl string,
	assetType string,
) (cai.Asset, error) {
	return cai.FetchIamPolicy(
		newUpdaterFunc,
		d,
		config,
		assetNameTmpl,
		assetType,
	)
}
