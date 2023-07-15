package google

import (
	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
)

func canonicalOrgPolicyConstraint(constraint string) string {
	return resourcemanager.CanonicalOrgPolicyConstraint(constraint)
}
