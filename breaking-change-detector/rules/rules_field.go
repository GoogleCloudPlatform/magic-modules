package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FieldRule provides structure for rules
// regarding field attribute changes
type FieldRule struct {
	name        string
	definition  string
	message     string
	identifier  string
	isRuleBreak func(old, new *schema.Schema) bool
}

// FieldRules is a list of FieldRule
// guarding against provider breaking changes
var FieldRules = []FieldRule{
	fieldRule_BecomingRequired,
}

var fieldRule_BecomingRequired = FieldRule{
	name:        "Field becoming Required Field",
	definition:  "A field cannot become required as existing terraform modules may not have this field defined. Thus breaking their modules in sequential plan or applies.",
	message:     "Field {{field}} changed from optional to required on {{resource}}",
	identifier:  "field-optional-to-required",
	isRuleBreak: fieldRule_BecomingRequired_func,
}

func fieldRule_BecomingRequired_func(old, new *schema.Schema) bool {
	if !old.Required && new.Required {
		return true
	}

	return false
}

func fieldRulesToRuleArray(frs []FieldRule) []Rule {
	var rules []Rule
	for _, fr := range frs {
		rules = append(rules, fr)
	}
	return rules
}

// Name - a human readable name for the rule
func (fr FieldRule) Name() string {
	return fr.name
}

// Definition - a definition for the rule
func (fr FieldRule) Definition() string {
	return fr.definition
}

// Identifier - a navigation oriented name for the rule
func (fr FieldRule) Identifier() string {
	return fr.identifier
}

// Message - a message to to inform the user
// of a breakage.
func (fr FieldRule) Message(version, resource, field string) string {
	msg := fr.message
	resource = fmt.Sprintf("`%s`", resource)
	field = fmt.Sprintf("`%s`", field)
	msg = strings.ReplaceAll(msg, "{{resource}}", resource)
	msg = strings.ReplaceAll(msg, "{{field}}", field)
	return msg + documentationReference(version, fr.identifier)
}

// IsRuleBreak - compares the fields and returns
// the if there was a rule breakage
func (fr FieldRule) IsRuleBreak(old, new *schema.Schema) bool {
	return fr.isRuleBreak(old, new)
}
