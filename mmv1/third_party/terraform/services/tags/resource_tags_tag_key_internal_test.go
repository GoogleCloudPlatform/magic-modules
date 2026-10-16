package tags

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func testTagsTagKeyImportResourceData(t *testing.T, id string) *schema.ResourceData {
	t.Helper()
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"purpose_data": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}, map[string]interface{}{})
	d.SetId(id)
	return d
}

func TestResourceTagsTagKeyImport(t *testing.T) {
	cases := map[string]struct {
		ImportId          string
		ExpectError       bool
		ExpectId          string
		ExpectName        string
		ExpectPurposeData map[string]interface{}
	}{
		"tagKeys prefixed name only": {
			ImportId:   "tagKeys/123456",
			ExpectId:   "tagKeys/123456",
			ExpectName: "123456",
		},
		"tagKeys prefixed name with GCE_FIREWALL purpose data": {
			ImportId:          "tagKeys/123456 network=my-project/default",
			ExpectId:          "tagKeys/123456",
			ExpectName:        "123456",
			ExpectPurposeData: map[string]interface{}{"network": "my-project/default"},
		},
		"tagKeys prefixed name with DATA_GOVERNANCE purpose data": {
			ImportId:          "tagKeys/123456 organization=auto",
			ExpectId:          "tagKeys/123456",
			ExpectName:        "123456",
			ExpectPurposeData: map[string]interface{}{"organization": "auto"},
		},
		"purpose data value containing an equals sign": {
			ImportId:          "tagKeys/123456 key=a=b=c",
			ExpectId:          "tagKeys/123456",
			ExpectName:        "123456",
			ExpectPurposeData: map[string]interface{}{"key": "a=b=c"},
		},
		"non numeric name is rejected": {
			ImportId:    "tagKeys/not-a-number",
			ExpectError: true,
		},
		"bare name without tagKeys prefix": {
			ImportId:   "123456",
			ExpectId:   "tagKeys/123456",
			ExpectName: "123456",
		},
		"bare name without tagKeys prefix with purpose data": {
			ImportId:          "123456 network=my-project/default",
			ExpectId:          "tagKeys/123456",
			ExpectName:        "123456",
			ExpectPurposeData: map[string]interface{}{"network": "my-project/default"},
		},
		"garbage input is rejected": {
			ImportId:    "not-a-valid-id",
			ExpectError: true,
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			d := testTagsTagKeyImportResourceData(t, tc.ImportId)
			results, err := resourceTagsTagKeyImport(d, nil)

			if tc.ExpectError {
				if err == nil {
					t.Fatalf("resourceTagsTagKeyImport(%q) = _, <nil>, want an error", tc.ImportId)
				}
				return
			}
			if err != nil {
				t.Fatalf("resourceTagsTagKeyImport(%q) = _, %v, want no error", tc.ImportId, err)
			}
			if len(results) != 1 {
				t.Fatalf("resourceTagsTagKeyImport(%q) returned %d ResourceData, want 1", tc.ImportId, len(results))
			}
			got := results[0]
			if got.Id() != tc.ExpectId {
				t.Errorf("Id() = %q, want %q", got.Id(), tc.ExpectId)
			}
			if name := got.Get("name").(string); name != tc.ExpectName {
				t.Errorf("Get(name) = %q, want %q", name, tc.ExpectName)
			}
			if tc.ExpectPurposeData != nil {
				if diff := got.Get("purpose_data").(map[string]interface{}); !mapsEqual(diff, tc.ExpectPurposeData) {
					t.Errorf("Get(purpose_data) = %v, want %v", diff, tc.ExpectPurposeData)
				}
			}
		})
	}
}

func mapsEqual(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			return false
		}
	}
	return true
}
