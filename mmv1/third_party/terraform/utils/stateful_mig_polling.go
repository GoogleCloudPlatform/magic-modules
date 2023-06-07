package google

import (
	tpgcompute "github.com/hashicorp/terraform-provider-google/google/services/compute"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Deprecated: For backward compatibility PollCheckInstanceConfigDeleted is still working,
// but all new code should use PollCheckInstanceConfigDeleted in the tpgcompute package instead.
func PollCheckInstanceConfigDeleted(resp map[string]interface{}, respErr error) transport_tpg.PollResult {
	return tpgcompute.PollCheckInstanceConfigDeleted(resp, respErr)
}
