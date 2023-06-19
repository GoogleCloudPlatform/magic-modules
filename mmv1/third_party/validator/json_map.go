package google

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
)

// JsonMap converts a given value to a map[string]interface{} that
// matches its JSON format.
func jsonMap(x interface{}) (map[string]interface{}, error) {
	return tpgresource.JsonMap(x)
}
