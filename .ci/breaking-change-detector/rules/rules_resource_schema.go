package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceSchemaRule provides structure for
// rules regarding resource attribute changes
type ResourceSchemaRule struct {
	name        string
	definition  string
	message     string
	identifier  string
	isRuleBreak func(old, new map[string]*schema.Schema) []string
}

// ResourceSchemaRule is a list of ResourceInventoryRule
// guarding against provider breaking changes
var ResourceSchemaRules = []ResourceSchemaRule{resourceSchemaRule_RemovingAField}

var resourceSchemaRule_RemovingAField = ResourceSchemaRule{
	name:        "Removing or Renaming an field",
	definition:  "In terraform fields should be retained whenever possible. A removable of an field will result in a configuration breakage wherever a dependency on that field exists. Renaming or Removing a field are functionally equivalent in terms of configuration breakages.",
	message:     "Field {{field}} within resource {{resource}} was either removed or renamed",
	identifier:  "resource-schema-field-removal-or-rename",
	isRuleBreak: resourceSchemaRule_RemovingAField_func,
}

func resourceSchemaRule_RemovingAField_func(old, new map[string]*schema.Schema) []string {
	keysNotPresent := []string{}
	for key := range old {
		_, exists := new[key]
		if !exists {
			keysNotPresent = append(keysNotPresent, key)
		}
	}
	return keysNotPresent
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
func (rs ResourceSchemaRule) Message(version, resource, field string) string {
	msg := rs.message
	resource = fmt.Sprintf("`%s`", resource)
	field = fmt.Sprintf("`%s`", field)
	msg = strings.ReplaceAll(msg, "{{resource}}", resource)
	msg = strings.ReplaceAll(msg, "{{field}}", field)
	return msg + documentationReference(version, rs.identifier)
}

// IsRuleBreak - compares the field entries and returns
// a list of fields violating the rule
func (rs ResourceSchemaRule) IsRuleBreak(old, new map[string]*schema.Schema) []string {
	return rs.isRuleBreak(old, new)
}
