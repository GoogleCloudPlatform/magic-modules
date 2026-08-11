package cloudrunv2

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// cloudRunV2ServiceResourceLimitsDiffSuppress suppresses diffs for the `cpu`
// key in the container resources.limits map when it was added by the API but
// not specified in the user's config. When only `memory` is set, the Cloud Run
// API automatically populates a default `cpu` value. Without this suppression
// the provider would produce a perpetual diff trying to clear that API-set key.
//
// Only the `cpu` key is suppressed because it is the only key the API
// auto-populates. Other keys such as `nvidia.com/gpu` must never be suppressed
// so that intentional removals (e.g. switching from GPU to non-GPU) are applied.
func cloudRunV2ServiceResourceLimitsDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// Only suppress the auto-populated cpu key, not user-set keys like nvidia.com/gpu.
	if !strings.HasSuffix(k, "limits.cpu") {
		return false
	}
	// Suppress when the user did not set cpu (new == "") but the API populated it (old != "").
	return new == "" && old != ""
}
