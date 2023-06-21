package google

import (
	"github.com/hashicorp/terraform-provider-google/google/services/billing"
)

func canonicalBillingAccountName(ba string) string {
	return billing.CanonicalBillingAccountName(ba)
}
