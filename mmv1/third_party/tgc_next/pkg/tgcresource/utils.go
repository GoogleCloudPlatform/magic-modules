package tgcresource

import (
	"fmt"
)

func GetComputeSelfLink(raw interface{}) interface{} {
	if raw == nil {
		return nil
	}

	v := raw.(string)
	if v != "" {
		return fmt.Sprintf("https://www.googleapis.com/compute/v1/%s", v)
	}

	return ""
}
