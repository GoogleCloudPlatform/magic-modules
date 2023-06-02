package google

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgiamresource"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
)

// ExpandIamPolicyBindings is used in google_<type>_iam_policy resources.
func expandIamPolicyBindings(d tpgresource.TerraformResourceData) ([]tpgresource.IAMBinding, error) {
	return tpgiamresource.ExpandIamPolicyBindings(d)
}

// ExpandIamRoleBindings is used in google_<type>_iam_binding resources.
func expandIamRoleBindings(d tpgresource.TerraformResourceData) ([]tpgresource.IAMBinding, error) {
	return tpgiamresource.ExpandIamRoleBindings(d)
}

// ExpandIamMemberBindings is used in google_<type>_iam_member resources.
func expandIamMemberBindings(d tpgresource.TerraformResourceData) ([]tpgresource.IAMBinding, error) {
	return tpgiamresource.ExpandIamMemberBindings(d)
}

// MergeIamAssets merges an existing asset with the IAM bindings of an incoming
// tpgresource.Asset.
func mergeIamAssets(
	existing, incoming tpgresource.Asset,
	MergeBindings func(existing, incoming []tpgresource.IAMBinding) []tpgresource.IAMBinding,
) tpgresource.Asset {
	return tpgiamresource.MergeIamAssets(
		existing, incoming,
		MergeBindings,
	)
}

// incoming is the last known state of an asset prior to deletion
func mergeDeleteIamAssets(
	existing, incoming tpgresource.Asset,
	MergeBindings func(existing, incoming []tpgresource.IAMBinding) []tpgresource.IAMBinding,
) tpgresource.Asset {
	return tpgiamresource.MergeDeleteIamAssets(
		existing, incoming,
		MergeBindings,
	)
}

// MergeAdditiveBindings adds members to bindings with the same roles and adds new
// bindings for roles that dont exist.
func mergeAdditiveBindings(existing, incoming []tpgresource.IAMBinding) []tpgresource.IAMBinding {
	return tpgiamresource.MergeAdditiveBindings(existing, incoming)
}

// MergeDeleteAdditiveBindings eliminates listed members from roles in the
// existing list. incoming is the last known state of the bindings being deleted.
func mergeDeleteAdditiveBindings(existing, incoming []tpgresource.IAMBinding) []tpgresource.IAMBinding {
	return tpgiamresource.MergeDeleteAdditiveBindings(existing, incoming)
}

// MergeAuthoritativeBindings clobbers members to bindings with the same roles
// and adds new bindings for roles that dont exist.
func mergeAuthoritativeBindings(existing, incoming []tpgresource.IAMBinding) []tpgresource.IAMBinding {
	return tpgiamresource.MergeAuthoritativeBindings(existing, incoming)
}

// MergeDeleteAuthoritativeBindings eliminates any bindings with matching roles
// in the existing list. incoming is the last known state of the bindings being
// deleted.
func mergeDeleteAuthoritativeBindings(existing, incoming []tpgresource.IAMBinding) []tpgresource.IAMBinding {
	return tpgiamresource.MergeDeleteAuthoritativeBindings(existing, incoming)
}

func fetchIamPolicy(
	newUpdaterFunc tpgiamresource.NewResourceIamUpdaterFunc,
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	assetNameTmpl string,
	assetType string,
) (tpgresource.Asset, error) {
	return tpgiamresource.FetchIamPolicy(
		newUpdaterFunc,
		d,
		config,
		assetNameTmpl,
		assetType,
	)
}
