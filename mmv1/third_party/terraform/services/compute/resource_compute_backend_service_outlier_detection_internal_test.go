package compute

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestUnitComputeBackendServiceOutlierDetectionEnforcement(t *testing.T) {
	resources := map[string]struct {
		resource *schema.Resource
		expand   func(interface{}, tpgresource.TerraformResourceData, *transport_tpg.Config) (interface{}, error)
	}{
		"global":   {ResourceComputeBackendService(), expandComputeBackendServiceOutlierDetection},
		"regional": {ResourceComputeRegionBackendService(), expandComputeRegionBackendServiceOutlierDetection},
	}
	for resourceName, resource := range resources {
		for _, test := range []struct {
			name             string
			configured       bool
			value            int
			withoutRawConfig bool
			wantPresent      bool
		}{
			{name: "omitted"},
			{name: "disabled", configured: true, value: 0, wantPresent: true},
			{name: "partial", configured: true, value: 50, wantPresent: true},
			{name: "unavailable raw config", configured: true, value: 0, withoutRawConfig: true},
		} {
			t.Run(resourceName+"/"+test.name, func(t *testing.T) {
				outlier := map[string]interface{}{"max_ejection_percent": 20}
				rawValue := cty.NullVal(cty.Number)
				if test.configured {
					outlier["enforcing_success_rate"] = test.value
					rawValue = cty.NumberIntVal(int64(test.value))
				}
				data := schema.TestResourceDataRaw(t, resource.resource.Schema, map[string]interface{}{
					"outlier_detection": []interface{}{outlier},
				})
				// The SDK test helper uses the legacy diff path, which does not retain
				// RawConfig. Supply the protocol's config separately from the state.
				data.SetId("test-backend")
				state := data.State()
				if !test.withoutRawConfig {
					state.RawConfig = cty.ObjectVal(map[string]cty.Value{
						"outlier_detection": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
							"enforcing_success_rate": rawValue,
						})}),
					})
				}
				data = resource.resource.Data(state)
				expanded, err := resource.expand(data.Get("outlier_detection"), data, &transport_tpg.Config{})
				if err != nil {
					t.Fatal(err)
				}
				body, err := json.Marshal(expanded)
				if err != nil {
					t.Fatal(err)
				}
				var request map[string]interface{}
				if err := json.Unmarshal(body, &request); err != nil {
					t.Fatal(err)
				}
				value, present := request["enforcingSuccessRate"]
				if present != test.wantPresent || (present && value != float64(test.value)) {
					t.Fatalf("enforcingSuccessRate: got value %v, present %v; want value %d, present %v; JSON: %s", value, present, test.value, test.wantPresent, body)
				}
			})
		}
	}
}
