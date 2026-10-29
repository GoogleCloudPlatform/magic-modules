package api

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/resource"
	"gopkg.in/yaml.v3"
)

func TestProductLowestVersion(t *testing.T) {
	t.Parallel()

	cases := []struct {
		description string
		obj         Product
		expected    string
	}{
		{
			description: "lowest version is ga",
			obj: Product{
				Versions: []*product.Version{
					&product.Version{
						Name:    "beta",
						BaseUrl: "beta_url",
					},
					&product.Version{
						Name:    "ga",
						BaseUrl: "ga_url",
					},
					&product.Version{
						Name:    "alpha",
						BaseUrl: "alpha_url",
					},
				},
			},
			expected: "ga",
		},
		{
			description: "lowest version is ga",
			obj: Product{
				Versions: []*product.Version{
					&product.Version{
						Name:    "beta",
						BaseUrl: "beta_url",
					},
					&product.Version{
						Name:    "alpha",
						BaseUrl: "alpha_url",
					},
				},
			},
			expected: "beta",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			versionObj := tc.obj.lowestVersion()

			if got, want := versionObj.Name, tc.expected; !reflect.DeepEqual(got, want) {
				t.Errorf("expected %v to be %v", got, want)
			}
		})
	}
}

func TestProductVersionObjOrClosest(t *testing.T) {
	t.Parallel()

	cases := []struct {
		description string
		obj         Product
		input       string
		expected    string
	}{
		{
			description: "closest version object to ga",
			obj: Product{
				Versions: []*product.Version{
					&product.Version{
						Name:    "beta",
						BaseUrl: "beta_url",
					},
					&product.Version{
						Name:    "ga",
						BaseUrl: "ga_url",
					},
				},
			},
			input:    "ga",
			expected: "ga",
		},
		{
			description: "closest version object to alpha",
			obj: Product{
				Versions: []*product.Version{
					&product.Version{
						Name:    "beta",
						BaseUrl: "beta_url",
					},
					&product.Version{
						Name:    "ga",
						BaseUrl: "ga_url",
					},
				},
			},
			input:    "alpha",
			expected: "beta",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			versionObj := tc.obj.VersionObjOrClosest(tc.input)

			if got, want := versionObj.Name, tc.expected; !reflect.DeepEqual(got, want) {
				t.Errorf("expected %v to be %v", got, want)
			}
		})
	}
}

func TestProductServiceName(t *testing.T) {
	t.Parallel()

	cases := []struct {
		description string
		obj         Product
		expected    string
	}{
		{
			description: "standard BaseUrl",
			obj: Product{
				Version: &product.Version{
					BaseUrl: "https://abc.googleapis.com/beta/",
				},
			},
			expected: "abc.googleapis.com",
		},
		{
			description: "BaseUrl with locational subdomain",
			obj: Product{
				Version: &product.Version{
					BaseUrl: "https://{{location}}-abc.googleapis.com/ga/",
				},
			},
			expected: "abc.googleapis.com",
		},
		{
			description: "BaseUrl and CaiBaseUrl",
			obj: Product{
				Version: &product.Version{
					BaseUrl:    "https://abc.googleapis.com/ga/",
					CaiBaseUrl: "https://def.googleapis.com/ga/",
				},
			},
			expected: "def.googleapis.com",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			if got, want := tc.obj.ServiceName(), tc.expected; got != want {
				t.Errorf("expected %v to be %v", got, want)
			}
		})
	}
}

func TestProductServiceVersion(t *testing.T) {
	t.Parallel()

	cases := []struct {
		description string
		obj         Product
		expected    string
	}{
		{
			description: "standard BaseUrl",
			obj: Product{
				Version: &product.Version{
					BaseUrl: "https://abc.googleapis.com/v1/",
				},
			},
			expected: "v1",
		},
		{
			description: "BaseUrl without trailing /",
			obj: Product{
				Version: &product.Version{
					BaseUrl: "https://abc.googleapis.com/v1",
				},
			},
			expected: "v1",
		},
		{
			description: "BaseUrl with version of beta",
			obj: Product{
				Version: &product.Version{
					BaseUrl: "https://abc.googleapis.com/beta/",
				},
			},
			expected: "beta",
		},
		{
			description: "BaseUrl without valid version",
			obj: Product{
				Version: &product.Version{
					BaseUrl: "https://abc.googleapis.com/other/",
				},
			},
			expected: "",
		},
		{
			description: "BaseUrl with additional value in path",
			obj: Product{
				Version: &product.Version{
					BaseUrl: "https://abc.googleapis.com/compute/v1/",
				},
			},
			expected: "v1",
		},
		{
			description: "standard BaseUrl",
			obj: Product{
				Version: &product.Version{
					BaseUrl:    "https://{{location}}-abc.googleapis.com/",
					CaiBaseUrl: "https://abc.googleapis.com/v1/",
				},
			},
			expected: "v1",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			if got, want := tc.obj.ServiceVersion(), tc.expected; got != want {
				t.Errorf("expected %v to be %v", got, want)
			}
		})
	}
}

func TestCustomUnmarshalYAMLRejectsUnknownFields(t *testing.T) {
	t.Parallel()

	strictDecode := func(raw string, out interface{}) error {
		dec := yaml.NewDecoder(bytes.NewReader([]byte(raw)))
		dec.KnownFields(true)
		return dec.Decode(out)
	}

	cases := []struct {
		name      string
		raw       string
		newTarget func() interface{}
		wantField string
		wantLine  string
	}{
		{
			name:      "Product",
			raw:       "name: Example\npackagepath: \"\"\n",
			newTarget: func() interface{} { return &Product{} },
			wantField: "packagepath",
			wantLine:  "line 2:",
		},
		{
			name:      "Async nested Operation",
			raw:       "operation:\n  base_url: '{{op_id}}'\n  result:\n    resource_inside_response: false\n",
			newTarget: func() interface{} { return &Async{} },
			wantField: "result",
			wantLine:  "line 3:",
		},
		{
			name:      "IamPolicy",
			raw:       "parent_resource_attribute: id\nmethod_name_suffix: Foo\n",
			newTarget: func() interface{} { return &resource.IamPolicy{} },
			wantField: "method_name_suffix",
			wantLine:  "line 2:",
		},
		{
			name:      "Examples",
			raw:       "name: basic\nnonexistent_example_key: true\n",
			newTarget: func() interface{} { return &resource.Examples{} },
			wantField: "nonexistent_example_key",
			wantLine:  "line 2:",
		},
		{
			name:      "Step",
			raw:       "name: basic\nvars:\n  foo: null\nnonexistent_step_key: true\n",
			newTarget: func() interface{} { return &resource.Step{} },
			wantField: "nonexistent_step_key",
			wantLine:  "line 4:",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := strictDecode(tc.raw, tc.newTarget())
			if err == nil {
				t.Fatalf("expected unknown field %q error on %s, got nil", tc.wantField, tc.name)
			}
			if !strings.Contains(err.Error(), tc.wantField) || !strings.Contains(err.Error(), tc.wantLine) {
				t.Errorf("expected error mentioning %q and %q, got: %v", tc.wantField, tc.wantLine, err)
			}
			if err := yaml.Unmarshal([]byte(tc.raw), tc.newTarget()); err != nil {
				t.Errorf("expected non-strict yaml.Unmarshal to ignore unknown field %q, got: %v", tc.wantField, err)
			}
		})
	}
}
