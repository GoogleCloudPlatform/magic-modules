package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceInventoryRule provides
// structure for rules regarding resource
// inventory changes
type ResourceInventoryRule struct {
	name        string
	definition  string
	message     string
	identifier  string
	isRuleBreak func(old, new map[string]*schema.Resource) []string
}

// ResourceInventoryRules is a list of ResourceInventoryRule
// guarding against provider breaking changes
var ResourceInventoryRules = []ResourceInventoryRule{resourceInventoryRule_RemovingAResource}

var resourceInventoryRule_RemovingAResource = ResourceInventoryRule{
	name:        "Removing or Renaming an Resource",
	definition:  "In terraform resources should be retained whenever possible. A removable of an resource will result in a configuration breakage wherever a dependency on that resource exists. Renaming or Removing a resources are functionally equivalent in terms of configuration breakages.",
	message:     "Resource {{resource}} was either removed or renamed",
	identifier:  "resource-map-resource-removal-or-rename",
	isRuleBreak: resourceInventoryRule_RemovingAField_func,
}

func resourceInventoryRule_RemovingAField_func(old, new map[string]*schema.Resource) []string {
	keysNotPresent := []string{}
	for key := range old {
		_, exists := new[key]
		if !exists {
			keysNotPresent = append(keysNotPresent, key)
		}
	}
	return keysNotPresent
}

func resourceInventoryRulesToRuleArray(rms []ResourceInventoryRule) []Rule {
	var rules []Rule
	for _, rm := range rms {
		rules = append(rules, rm)
	}
	return rules
}

// Name - a human readable name for the rule
func (rm ResourceInventoryRule) Name() string {
	return rm.name
}

// Definition - a definition for the rule
func (rm ResourceInventoryRule) Definition() string {
	return rm.definition
}

// Identifier - a navigation oriented name for the rule
func (rm ResourceInventoryRule) Identifier() string {
	return rm.identifier
}

// Message - a message to to inform the user
// of a breakage.
func (rm ResourceInventoryRule) Message(resource string) string {
	msg := rm.message
	resource = fmt.Sprintf("`%s`", resource)
	msg = strings.ReplaceAll(msg, "{{resource}}", resource)
	return msg + documentationReference(rm.identifier)
}

// IsRuleBreak - compares resource entries and returns
// a list of resources violating the rule
func (rm ResourceInventoryRule) IsRuleBreak(old, new map[string]*schema.Resource) []string {
	if rm.isRuleBreak == nil {
		return []string{}
	}
	return rm.isRuleBreak(old, new)
}

// Undetectable - informs if there are functions in place
// to detect this rule.
func (rm ResourceInventoryRule) Undetectable() bool {
	return rm.isRuleBreak == nil
}
