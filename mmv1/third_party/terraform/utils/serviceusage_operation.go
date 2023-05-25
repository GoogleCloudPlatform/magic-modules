package google

import (
	"encoding/json"
	"time"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/serviceusage/v1"
)

func serviceUsageOperationWait(config *transport_tpg.Config, op *serviceusage.Operation, project, activity, userAgent string, timeout time.Duration) error {
	// maintained for compatibility with old code that was written before the
	// autogenerated waiters.
	b, err := op.MarshalJSON()
	if err != nil {
		return err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	return ServiceUsageOperationWaitTime(config, m, project, activity, userAgent, timeout)
}
