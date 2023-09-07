package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ProviderConfigRule provides
// structure for rules regarding resource
// inventory changes
type ProviderConfigRule struct {
	name        string
	definition  string
	message     string
	identifier  string
	isRuleBreak func(old, new map[string]*schema.Resource) []string
}

// ProviderConfigRules is a list of ProviderConfigRule
// guarding against provider breaking changes
var ProviderConfigRules = []ProviderConfigRule{providerConfigRule_ConfigurationChanges}

var providerConfigRule_ConfigurationChanges = ProviderConfigRule{
	name:        "Changing fundamental provider behavior",
	definition:  "Including, but not limited to modification of: authentication, environment variable usage, and constricting retry behavior.",
	identifier:  "provider-config-fundamental",
	isRuleBreak: nil,
}

func providerConfigRulesToRuleArray(pcrs []ProviderConfigRule) []Rule {
	var rules []Rule
	for _, prc := range pcrs {
		rules = append(rules, prc)
	}
	return rules
}

// Name - a human readable name for the rule
func (prc ProviderConfigRule) Name() string {
	return prc.name
}

// Definition - a definition for the rule
func (prc ProviderConfigRule) Definition() string {
	return prc.definition
}

// Identifier - a navigation oriented name for the rule
func (prc ProviderConfigRule) Identifier() string {
	return prc.identifier
}

// Message - a message to to inform the user
// of a breakage.
func (prc ProviderConfigRule) Message(resource string) string {
	msg := prc.message
	resource = fmt.Sprintf("`%s`", resource)
	msg = strings.ReplaceAll(msg, "{{resource}}", resource)
	return msg + documentationReference(prc.identifier)
}

// IsRuleBreak - compares resource entries and returns
// a list of resources violating the rule
func (prc ProviderConfigRule) IsRuleBreak(old, new map[string]*schema.Resource) []string {
	if prc.isRuleBreak == nil {
		return []string{}
	}
	return prc.isRuleBreak(old, new)
}

func (prc ProviderConfigRule) Undetectable() bool {
	return prc.isRuleBreak == nil
}
