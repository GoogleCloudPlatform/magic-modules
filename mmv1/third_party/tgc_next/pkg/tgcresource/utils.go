package tgcresource

import (
	"fmt"
	"strings"

	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/transport"
)

// Remove the Terraform attribution label "goog-terraform-provisioned" from labels
func RemoveTerraformAttributionLabel(raw interface{}) interface{} {
	if raw == nil {
		return nil
	}

	if labels, ok := raw.(map[string]string); ok {
		delete(labels, "goog-terraform-provisioned")
		return labels
	}

	if labels, ok := raw.(map[string]interface{}); ok {
		delete(labels, "goog-terraform-provisioned")
		return labels
	}

	return nil
}

// Gets the full url from relative url
func GetFullUrl(config *transport_tpg.Config, raw interface{}, baseUrl string) interface{} {
	if raw == nil || baseUrl == "" {
		return raw
	}

	v := raw.(string)
	if v != "" && !strings.HasPrefix(v, "https://") {
		if config.UniverseDomain == "" || config.UniverseDomain == "googleapis.com" {
			return fmt.Sprintf("%s%s", baseUrl, v)
		}
	}

	return v
}
