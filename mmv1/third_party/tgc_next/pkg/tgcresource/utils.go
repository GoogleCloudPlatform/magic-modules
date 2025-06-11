package tgcresource

import (
	"fmt"
	"strings"

	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/transport"
)

func GetComputeSelfLink(config *transport_tpg.Config, raw interface{}) interface{} {
	if raw == nil {
		return nil
	}

	v := raw.(string)
	if v != "" && !strings.HasPrefix(v, "https://") {
		if config.UniverseDomain == "" || config.UniverseDomain == "googleapis.com" {
			return fmt.Sprintf("https://www.googleapis.com/compute/v1/%s", v)
		}
	}

	return v
}
