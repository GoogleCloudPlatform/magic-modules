package detector

import (
	"strings"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func initialModifyStatus(field *schema.Schema, resourceDiff diff.ResourceDiff) TransitionResult {
	if resourceDiff.NoUpdate {
		return TransitionResult{Status: StatusExempt, ExemptCode: "resource_not_updatable"}
	}
	if field.ForceNew {
		return TransitionResult{Status: StatusExempt, ExemptCode: "force_new"}
	}
	return TransitionResult{Status: StatusMissing}
}

func initialUnsetStatus(fieldPath string, field *schema.Schema) TransitionResult {
	if field.Required {
		if !strings.Contains(fieldPath, ".") {
			return TransitionResult{Status: StatusExempt, ExemptCode: "required"}
		}
		return TransitionResult{Status: StatusExempt, ExemptCode: "required_in_parent"}
	}
	if field.Optional && field.Computed {
		return TransitionResult{Status: StatusExempt, ExemptCode: "optional_computed"}
	}
	return TransitionResult{Status: StatusMissing}
}

// usesOmissionUnset returns true when a field or block cannot be unset via an in-place
// update step (e.g. NoUpdate resource or ForceNew field) and therefore satisfies UNSET
// by being omitted in at least one valid apply step where its parent is present.
func usesOmissionUnset(field *schema.Schema, resourceDiff diff.ResourceDiff) bool {
	if resourceDiff.NoUpdate {
		return true
	}
	if field != nil && field.ForceNew {
		return true
	}
	return false
}
