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

// Terraform must set the top level schema field, but since this object contains collapsed properties
// it's difficult to know what the top level should be. Instead we just loop over the map returned from flatten.
func MergeFlattenedProperties(hclData map[string]interface{}, flattenedProp interface{}) error {
	if flattenedProp == nil {
		return nil
	}
	flattenedPropSlice, ok := flattenedProp.([]interface{})
	if !ok || len(flattenedPropSlice) == 0 {
		return fmt.Errorf("unexpected type returned from flattener: %T", flattenedProp)
	}
	flattedPropMap, ok := flattenedPropSlice[0].(map[string]interface{})
	if !ok || len(flattedPropMap) == 0 {
		return fmt.Errorf("unexpected type returned from flattener: %T", flattenedPropSlice)
	}
	for k, v := range flattedPropMap {
		hclData[k] = v
	}
	return nil
}

// Checks if all values in the map are nil
func AllValuesAreNil(m map[string]interface{}) bool {
	if len(m) == 0 {
		return true
	}

	for _, v := range m {
		if v != nil {
			return false
		}
	}

	return true
}
