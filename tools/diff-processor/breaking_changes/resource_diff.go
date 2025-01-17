package breaking_changes

import (
	"fmt"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
)

// ResourceDiffRule is a rule that operates on an entire ResourceDiff
type ResourceDiffRule struct {
	Identifier string
	Messages   func(resource string, resourceDiff diff.ResourceDiff) []string
}

// ResourceDiffRules is a list of all ResourceDiff rules
var ResourceDiffRules = []ResourceDiffRule{RemovingAField}

var RemovingAField = ResourceDiffRule{
	Identifier: "resource-schema-field-removal-or-rename",
	Messages:   RemovingAFieldMessages,
}

// TODO: Make field removal a FieldDiffRule b/300124253
func RemovingAFieldMessages(resource string, resourceDiff diff.ResourceDiff) []string {
	fieldsRemoved := []string{}
	for field, fieldDiff := range resourceDiff.Fields {
		if fieldDiff.Old != nil && fieldDiff.New == nil {
			fieldsRemoved = append(fieldsRemoved, field)
		}
	}

	tmpl := "Field `%s` within resource `%s` was either removed or renamed"
	var messages []string
	for _, field := range fieldsRemoved {
		messages = append(messages, fmt.Sprintf(tmpl, field, resource))
	}
	return messages
}
