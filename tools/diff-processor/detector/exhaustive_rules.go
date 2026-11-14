package detector

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func initialModifyStatus(field *schema.Schema) TransitionResult {
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
