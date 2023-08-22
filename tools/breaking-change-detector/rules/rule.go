package rules

// Rule is an interface that all Rule
// types implement. This give consistency
// for the potential documentation generation.
type Rule interface {
	Name() string
	Definition() string
	Identifier() string
	Undetectable() bool
}

// RuleCategory holds documentation and
// inventory for each of the rule types
type RuleCategory struct {
	Name       string
	Definition string
	Rules      []Rule
}

// Rules is a type to hold all existing
// rules. Exposed for potential documentation generator
type Rules struct {
	Categories []RuleCategory
}

// GetRules returns a list of all rules for
// potential documentation generation
func GetRules() *Rules {
	categories := []RuleCategory{
		{
			Name:       "Provider Configuration Level Breakages",
			Definition: "Top level behavior such as provider configuration and authentication changes.",
			Rules:      providerConfigRulesToRuleArray(ProviderConfigRules),
		},
		{
			Name:       "Resource List Level Breakages",
			Definition: "Resource/datasource naming conventions and entry differences.",
			Rules:      resourceInventoryRulesToRuleArray(ResourceInventoryRules),
		},
		{
			Name:       "Resource Level Breakages",
			Definition: "Individual resource breakages like field entry removals or behavior within a resource.",
			Rules:      resourceSchemaRulesToRuleArray(ResourceSchemaRules),
		},
		{
			Name:       "Field Level Breakages",
			Definition: "Field level conventions like attribute changes and naming conventions.",
			Rules:      fieldRulesToRuleArray(FieldRules),
		},
	}

	return &Rules{
		Categories: categories,
	}
}
