package firebaserules

import (
	dcl "github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
)

// EncodeReleaseUpdateRequest encapsulates fields in a release {} block, as expected
// by https://firebase.google.com/docs/reference/rules/rest/v1/projects.releases/patch
func EncodeReleaseUpdateRequest(m map[string]interface{}) map[string]interface{} {
	req := make(map[string]interface{})
	dcl.PutMapEntry(req, []string{"release"}, m)
	return req
}
