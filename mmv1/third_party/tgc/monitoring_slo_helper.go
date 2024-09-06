package monitoring

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

func expandMonitoringSloRollingPeriodDays(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	i, ok := v.(int)
	if !ok {
		return nil, fmt.Errorf("unexpected value is not int: %v", v)
	}
	if i == 0 {
		return "", nil
	}
	// Day = 86400s
	return fmt.Sprintf("%ds", i*86400), nil
}
