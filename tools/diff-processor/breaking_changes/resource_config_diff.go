package breaking_changes

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceInventoryRule provides
// structure for rules regarding resource
// inventory changes
type ResourceConfigDiffRule struct {
	Identifier string
	// TODO: Make this take a ResourceConfigDiff instead of old, new.
	Messages func(resource string, old, new *schema.Resource) []string
}

// ResourceInventoryRules is a list of ResourceInventoryRule
// guarding against provider breaking changes
var ResourceConfigDiffRules = []ResourceConfigDiffRule{ResourceConfigRemovingAResource}

var ResourceConfigRemovingAResource = ResourceConfigDiffRule{
	Identifier: "resource-map-resource-removal-or-rename",
	Messages:   ResourceConfigRemovingAResourceMessages,
}

func ResourceConfigRemovingAResourceMessages(resource string, old, new *schema.Resource) []string {
	if new == nil && old != nil {
		tmpl := "Resource `%s` was either removed or renamed"
		return []string{fmt.Sprintf(tmpl, resource)}
	}
	return nil
}
