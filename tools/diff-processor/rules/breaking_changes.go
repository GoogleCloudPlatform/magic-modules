package rules

import (
	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
)

func ComputeBreakingChanges(schemaDiff diff.SchemaDiff) []string {
	messages := []string{}
	for resource, resourceDiff := range schemaDiff {
		for _, rule := range ResourceInventoryRules {
			if rule.isRuleBreak(resourceDiff.ResourceConfig.Old, resourceDiff.ResourceConfig.New) {
				messages = append(messages, rule.Message(resource))
			}
		}

		// If the resource was added or removed, don't check field schema diffs.
		if resourceDiff.ResourceConfig.Old == nil || resourceDiff.ResourceConfig.New == nil {
			continue
		}

		// TODO: Move field removal to field_rules and merge resource schema / resource inventory rules
		// b/300124253
		for _, rule := range ResourceSchemaRules {
			violatingFields := rule.IsRuleBreak(resourceDiff)
			if len(violatingFields) > 0 {
				for _, field := range violatingFields {
					newMessage := rule.Message(resource, field)
					messages = append(messages, newMessage)
				}
			}
		}

		for field, fieldDiff := range resourceDiff.Fields {
			for _, rule := range FieldRules {
				// TODO: refactor rules to use interface-based implementation that separates checking whether
				// a rule broke from composing a message for a rule break.
				breakageMessage := rule.IsRuleBreak(
					fieldDiff.Old,
					fieldDiff.New,
					MessageContext{
						Resource: resource,
						Field:    field,
					},
				)
				if breakageMessage != "" {
					messages = append(messages, breakageMessage)
				}
			}
		}
	}
	return messages
}
