package breaking_changes

import (
	"fmt"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
)

// ResourceConfigDiffRule provides
// structure for rules regarding resource config changes
type ResourceConfigDiffRule struct {
	Identifier string
	Messages func(resource string, diff diff.ResourceConfigDiff) []string
}

// ResourceConfigDiffRules is a list of ResourceConfigDiffRule
// guarding against provider breaking changes
var ResourceConfigDiffRules = []ResourceConfigDiffRule{ResourceConfigRemovingAResource}

var ResourceConfigRemovingAResource = ResourceConfigDiffRule{
	Identifier: "resource-map-resource-removal-or-rename",
	Messages:   ResourceConfigRemovingAResourceMessages,
}

func ResourceConfigRemovingAResourceMessages(resource string, diff diff.ResourceConfigDiff) []string {
	if diff.New == nil && diff.Old != nil {
		tmpl := "Resource `%s` was either removed or renamed"
		return []string{fmt.Sprintf(tmpl, resource)}
	}
	return nil
}
