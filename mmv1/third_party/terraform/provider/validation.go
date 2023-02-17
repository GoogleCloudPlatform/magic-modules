package provider

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func validateNonNegativeDuration() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		dur, err := time.ParseDuration(v)
		if err != nil {
			es = append(es, fmt.Errorf("expected %s to be a duration, but parsing gave an error: %s", k, err.Error()))
			return
		}

		if dur < 0 {
			es = append(es, fmt.Errorf("duration %v must be a non-negative duration", dur))
			return
		}

		return
	}
}
