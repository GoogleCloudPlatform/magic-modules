package cai

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

func TestAssetName(t *testing.T) {
	cases := []struct {
		name            string
		template        string
		expectedPattern string
		data            tpgresource.TerraformResourceData
	}{
		{
			name:            "PresentValues",
			template:        "//{{a}}/{{b}}",
			expectedPattern: "//value-a/value-b",
			data: &mockTerraformResourceData{
				m: map[string]interface{}{
					"a": "value-a",
					"b": "value-b",
				},
			},
		},
		{
			name:            "MissingValue",
			template:        "//{{a}}/{{b}}",
			expectedPattern: `//value-a/placeholder-\S{8}`,
			data: &mockTerraformResourceData{
				m: map[string]interface{}{
					"a": "value-a",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := regexp.MustCompile(c.expectedPattern)

			output, err := AssetName(c.data, &transport_tpg.Config{}, c.template)
			if err != nil {
				t.Fatal(err)
			}

			if !r.MatchString(output) {
				t.Fatalf("got %v, expected pattern %v", output, c.expectedPattern)
			}
		})
	}
}

func TestRandString(t *testing.T) {
	memory := make(map[string]bool)
	for i := 0; i < 100; i++ {
		s := RandString(i)
		if n := len(s); n != i {
			t.Fatalf("expected len = %v, got %v", i, n)
		}
		if memory[s] {
			t.Fatalf("already seen string: %v, probably not random!", s)
		}
		memory[s] = true
	}
}

type mockTerraformResourceData struct {
	m map[string]interface{}
	tpgresource.TerraformResourceData
}

func (d *mockTerraformResourceData) GetOkExists(k string) (interface{}, bool) {
	v, ok := d.m[k]
	return v, ok
}

func (d *mockTerraformResourceData) GetOk(k string) (interface{}, bool) {
	v, ok := d.m[k]
	return v, ok
}
