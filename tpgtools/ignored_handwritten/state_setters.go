package google

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func setStateForTime(d *schema.ResourceData, v time.Time, name string) error {
	if !v.IsZero() {
		return d.Set(name, fmt.Sprintf(v.Format(time.RFC3339)))
	} else {
		return d.Set(name, nil)
	}
}

func generateIfNotSet(d *schema.ResourceData, field, prefix string) (string, error) {
	if _, ok := d.GetOkExists(field); !ok {
		if prefix == "" {
			prefix = "tf-generated-"
		}
		v := resource.PrefixedUniqueId(prefix)
		if len(v) > 30 {
			v = v[:30]
		}

		if err := d.Set(field, v); err != nil {
			return "", err
		}
	}
	return d.Get(field).(string), nil
}
