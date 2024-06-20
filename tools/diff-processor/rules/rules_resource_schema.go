package rules

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/constants"
	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
)

// ResourceSchemaRule provides structure for
// rules regarding resource attribute changes
type ResourceSchemaRule struct {
	name        string
	definition  string
	message     string
	identifier  string
	isRuleBreak func(resourceDiff diff.ResourceDiff) []string
}

// ResourceSchemaRules is a list of ResourceInventoryRule
// guarding against provider breaking changes
var ResourceSchemaRules = []ResourceSchemaRule{resourceSchemaRule_RemovingAField, resourceSchemaRule_ChangingResourceIDFormat, resourceSchemaRule_ChangingImportIDFormat}

var resourceSchemaRule_ChangingResourceIDFormat = ResourceSchemaRule{
	name:       "Changing resource ID format",
	definition: "Terraform uses resource ID to read resource state from the api. Modification of the ID format will break the ability to parse the IDs from any deployments.",
	identifier: "resource-id",
}

var resourceSchemaRule_ChangingImportIDFormat = ResourceSchemaRule{
	name:       "Changing resource ID import format",
	definition: "Automation external to our provider may rely on importing resources with a certain format. Removal or modification of existing formats will break this automation.",
	identifier: "resource-import-format",
}

var resourceSchemaRule_RemovingAField = ResourceSchemaRule{
	name:        "Removing or Renaming an field",
	definition:  "In terraform fields should be retained whenever possible. A removable of an field will result in a configuration breakage wherever a dependency on that field exists. Renaming or Removing a field are functionally equivalent in terms of configuration breakages.",
	message:     "Field {{field}} within resource {{resource}} was either removed or renamed",
	identifier:  "resource-schema-field-removal-or-rename",
	isRuleBreak: resourceSchemaRule_RemovingAField_func,
}

func resourceSchemaRule_RemovingAField_func(resourceDiff diff.ResourceDiff) []string {
	fieldsRemoved := []string{}
	for field, fieldDiff := range resourceDiff.Fields {
		if fieldDiff.Old != nil && fieldDiff.New == nil {
			fieldsRemoved = append(fieldsRemoved, field)
		}
	}
	return fieldsRemoved
}

func resourceSchemaRulesToRuleArray(rss []ResourceSchemaRule) []Rule {
	var rules []Rule
	for _, rs := range rss {
		rules = append(rules, rs)
	}
	return rules
}

// Name - a human readable name for the rule
func (rs ResourceSchemaRule) Name() string {
	return rs.name
}

// Definition - a definition for the rule
func (rs ResourceSchemaRule) Definition() string {
	return rs.definition
}

// Identifier - a navigation oriented name for the rule
func (rs ResourceSchemaRule) Identifier() string {
	return rs.identifier
}

// Message - a message to to inform the user
// of a breakage.
func (rs ResourceSchemaRule) Message(resource, field string) *BreakingChange {
	msg := rs.message
	msg = strings.ReplaceAll(msg, "{{resource}}", fmt.Sprintf("`%s`", resource))
	msg = strings.ReplaceAll(msg, "{{field}}", fmt.Sprintf("`%s`", field))
	return &BreakingChange{
		Resource:               resource,
		Field:                  field,
		Message:                msg,
		DocumentationReference: constants.GetFileUrl(rs.identifier),
		RuleTemplate:           rs.message,
		RuleDefinition:         rs.definition,
		RuleName:               rs.name,
	}
}

// IsRuleBreak - compares the field entries and returns
// a list of fields violating the rule
func (rs ResourceSchemaRule) IsRuleBreak(resourceDiff diff.ResourceDiff) []string {
	if rs.isRuleBreak == nil {
		return []string{}
	}
	return rs.isRuleBreak(resourceDiff)
}

// Undetectable - informs if there are functions in place
// to detect this rule.
func (rs ResourceSchemaRule) Undetectable() bool {
	return rs.isRuleBreak == nil
}
