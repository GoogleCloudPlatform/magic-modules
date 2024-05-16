package rules

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/constants"
)

// MessageContext - is an envelope for additional
// metadata about context around the rule breakage
type MessageContext struct {
	Resource   string
	Field      string
	definition string
	name       string
	identifier string
	message    string
}

func populateMessageContext(message string, mc MessageContext) *BreakingChange {
	template := message
	resource := fmt.Sprintf("`%s`", mc.Resource)
	field := fmt.Sprintf("`%s`", mc.Field)
	message = strings.ReplaceAll(message, "{{resource}}", resource)
	message = strings.ReplaceAll(message, "{{field}}", field)
	return &BreakingChange{
		Resource:               mc.Resource,
		Field:                  mc.Field,
		Message:                message,
		RuleTemplate:           template,
		DocumentationReference: constants.GetFileUrl(mc.identifier),
		RuleDefinition:         mc.definition,
		RuleName:               mc.name,
	}
}
