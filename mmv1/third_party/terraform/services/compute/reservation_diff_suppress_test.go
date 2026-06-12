package compute

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func setResourceDataConfig(d *schema.ResourceData, config *terraform.ResourceConfig) {
	val := reflect.ValueOf(d).Elem()
	configField := val.FieldByName("config")
	ptr := unsafe.Pointer(configField.UnsafeAddr())
	*(**terraform.ResourceConfig)(ptr) = config
}

func TestReservationDiffSuppress(t *testing.T) {
	t.Parallel()

	reservationSchema := ResourceComputeReservation().Schema

	cases := map[string]struct {
		rawConfig           cty.Value
		expectSuppressMap   bool
		expectSuppressProjs bool
	}{
		"only_project_map_configured": {
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"share_settings": cty.ListVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"share_type": cty.StringVal("SPECIFIC_PROJECTS"),
						"project_map": cty.SetVal([]cty.Value{
							cty.ObjectVal(map[string]cty.Value{
								"id":         cty.StringVal("my-project-id"),
								"project_id": cty.StringVal("my-project-id"),
							}),
						}),
					}),
				}),
			}),
			expectSuppressMap:   false,
			expectSuppressProjs: true,
		},
		"only_projects_configured": {
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"share_settings": cty.ListVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"share_type": cty.StringVal("SPECIFIC_PROJECTS"),
						"projects": cty.ListVal([]cty.Value{
							cty.StringVal("my-project-id"),
						}),
					}),
				}),
			}),
			expectSuppressMap:   true,
			expectSuppressProjs: false,
		},
		"both_configured": {
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"share_settings": cty.ListVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"share_type": cty.StringVal("SPECIFIC_PROJECTS"),
						"project_map": cty.SetVal([]cty.Value{
							cty.ObjectVal(map[string]cty.Value{
								"id":         cty.StringVal("my-project-id"),
								"project_id": cty.StringVal("my-project-id"),
							}),
						}),
						"projects": cty.ListVal([]cty.Value{
							cty.StringVal("another-project-id"),
						}),
					}),
				}),
			}),
			expectSuppressMap:   false,
			expectSuppressProjs: false,
		},
		"neither_configured": {
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"share_settings": cty.ListVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"share_type": cty.StringVal("LOCAL"),
					}),
				}),
			}),
			expectSuppressMap:   false,
			expectSuppressProjs: false,
		},
	}

	for tn, tc := range cases {
		d := schema.TestResourceDataRaw(t, reservationSchema, map[string]interface{}{})
		
		resourceConfig := &terraform.ResourceConfig{
			CtyValue: tc.rawConfig,
		}
		setResourceDataConfig(d, resourceConfig)

		suppressMap := computeReservationProjectMapDiffSuppress("", "", "", d)
		if suppressMap != tc.expectSuppressMap {
			t.Errorf("case %s failed: expected projectMap suppress to be %v, got %v", tn, tc.expectSuppressMap, suppressMap)
		}

		suppressProjs := computeReservationProjectsDiffSuppress("", "", "", d)
		if suppressProjs != tc.expectSuppressProjs {
			t.Errorf("case %s failed: expected projects suppress to be %v, got %v", tn, tc.expectSuppressProjs, suppressProjs)
		}
	}
}
