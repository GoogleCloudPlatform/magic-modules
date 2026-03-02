package tpgresource

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func normalizeForwardingRule(v string) string {
	if strings.HasPrefix(v, "https://") {
		if idx := strings.Index(v, "/projects/"); idx != -1 {
			v = v[idx:]
		}
	}

	v = strings.TrimPrefix(v, "/")

	return v
}

func CompareForwardingRuleSelfLinkOrName(_, old, new string, _ *schema.ResourceData) bool {
	if old == new {
		return true
	}

	oldNorm := normalizeForwardingRule(old)
	newNorm := normalizeForwardingRule(new)

	if strings.Contains(oldNorm, "/forwardingRules/") &&
		strings.Contains(newNorm, "/forwardingRules/") {
		return oldNorm == newNorm
	}

	oldName := oldNorm[strings.LastIndex(oldNorm, "/")+1:]
	newName := newNorm[strings.LastIndex(newNorm, "/")+1:]

	return oldName == newName
}
