package breaking_changes

import (
	"fmt"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
)

type BreakingChange struct {
	Resource               string
	Field                  string
	Message                string
	DocumentationReference string
	RuleName               string
}

const breakingChangesPath = "develop/breaking-changes/breaking-changes"

func NewBreakingChange(message, identifier string) BreakingChange {
	return BreakingChange{
		Message:                message,
		DocumentationReference: fmt.Sprintf("https://googlecloudplatform.github.io/magic-modules/%s#%s", breakingChangesPath, identifier),
	}
}

func ComputeBreakingChanges(schemaDiff diff.SchemaDiff) []BreakingChange {
	var breakingChanges []BreakingChange
	for resource, resourceDiff := range schemaDiff {
		for _, rule := range ResourceConfigDiffRules {
			for _, message := range rule.Messages(resource, resourceDiff.ResourceConfig.Old, resourceDiff.ResourceConfig.New) {
				breakingChanges = append(breakingChanges, NewBreakingChange(message, rule.Identifier))
			}
		}

		// If the resource was added or removed, don't check rules that include field information.
		if resourceDiff.ResourceConfig.Old == nil || resourceDiff.ResourceConfig.New == nil {
			continue
		}

		for _, rule := range ResourceDiffRules {
			for _, message := range rule.Messages(resource, resourceDiff) {
				breakingChanges = append(breakingChanges, NewBreakingChange(message, rule.Identifier))
			}
		}

		for field, fieldDiff := range resourceDiff.Fields {
			for _, rule := range FieldDiffRules {
				for _, message := range rule.Messages(resource, field, fieldDiff.Old, fieldDiff.New) {
					breakingChanges = append(breakingChanges, NewBreakingChange(message, rule.Identifier))
				}
			}
		}
	}
	return breakingChanges
}
