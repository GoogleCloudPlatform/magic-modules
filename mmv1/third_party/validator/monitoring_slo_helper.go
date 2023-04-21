package google

import (
	"fmt"

	transport_tpg "github.com/GoogleCloudPlatform/terraform-validator/converters/google/resources/transport"
)

func expandMonitoringSloRollingPeriodDays(v interface{}, d TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
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
