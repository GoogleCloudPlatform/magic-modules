package rules

import (
	"fmt"
	"strings"
)

// MessageContext - is an envelope for additional
// metadata about context around the rule breakage
type MessageContext struct {
	Resource   string
	Field      string
	identifier string
	message    string
}

func populateMessageContext(message string, mc MessageContext) string {
	resource := fmt.Sprintf("`%s`", mc.Resource)
	field := fmt.Sprintf("`%s`", mc.Field)
	message = strings.ReplaceAll(message, "{{resource}}", resource)
	message = strings.ReplaceAll(message, "{{field}}", field)
	return message + documentationReference(mc.identifier)
}
