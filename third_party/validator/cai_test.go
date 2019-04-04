package google

import (
	"regexp"
	"testing"

	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

func TestReplaceWithPlaceholder(t *testing.T) {
	cases := []struct {
		name            string
		template        string
		expectedPattern string
		data            TerraformResourceData
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

			output, err := replaceWithPlaceholder(c.data, nil, c.template)
			if err != nil {
				t.Fatal(err)
			}

			if !r.MatchString(output) {
				t.Fatalf("got %v, expected pattern %v", output, c.expectedPattern)
			}
		})
	}
}

type mockTerraformResourceData struct {
	m map[string]interface{}
	TerraformResourceData
}

func (d *mockTerraformResourceData) GetOk(k string) (interface{}, bool) {
	v, ok := d.m[k]
	return v, ok
}

func TestRandString(t *testing.T) {
	memory := make(map[string]bool)
	for i := 0; i < 100; i++ {
		s := randString(i)
		if n := len(s); n != i {
			t.Fatalf("expected len = %v, got %v", i, n)
		}
		if memory[s] {
			t.Fatalf("already seen string: %v, probably not random!", s)
		}
		memory[s] = true
	}
}

func TestAncestryPath(t *testing.T) {
	cases := []struct {
		name           string
		input          []*cloudresourcemanager.Ancestor
		expectedOutput string
	}{
		{
			name:           "Empty",
			input:          []*cloudresourcemanager.Ancestor{},
			expectedOutput: "",
		},
		{
			name: "ProjectOrganization",
			input: []*cloudresourcemanager.Ancestor{
				{
					ResourceId: &cloudresourcemanager.ResourceId{
						Id:   "my-prj",
						Type: "project",
					},
				},
				{
					ResourceId: &cloudresourcemanager.ResourceId{
						Id:   "my-org",
						Type: "organization",
					},
				},
			},
			expectedOutput: "organization/my-org/project/my-prj",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			output := ancestryPath(c.input)
			if output != c.expectedOutput {
				t.Errorf("expected output %q, got %q", c.expectedOutput, output)
			}
		})
	}
}
