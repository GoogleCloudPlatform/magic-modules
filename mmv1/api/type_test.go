package api

import (
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
)

func TestTypeMinVersionObj(t *testing.T) {
	t.Parallel()

	p := Product{
		Name: "test",
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
	}

	cases := []struct {
		description string
		obj         Type
		expected    string
	}{
		{
			description: "type minVersion is empty and resource minVersion is empty",
			obj: Type{
				Name:       "test",
				MinVersion: "",
				ResourceMetadata: &Resource{
					Name:            "test",
					MinVersion:      "",
					ProductMetadata: &p,
				},
			},
			expected: "ga",
		},
		{
			description: "type minVersion is empty and resource minVersion is beta",
			obj: Type{
				Name:       "test",
				MinVersion: "",
				ResourceMetadata: &Resource{
					Name:            "test",
					MinVersion:      "beta",
					ProductMetadata: &p,
				},
			},
			expected: "beta",
		},
		{
			description: "type minVersion is not empty",
			obj: Type{
				Name:       "test",
				MinVersion: "beta",
				ResourceMetadata: &Resource{
					Name:            "test",
					MinVersion:      "",
					ProductMetadata: &p,
				},
			},
			expected: "beta",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			versionObj := tc.obj.MinVersionObj()

			if got, want := versionObj.Name, tc.expected; !reflect.DeepEqual(got, want) {
				t.Errorf("expected %v to be %v", got, want)
			}
		})
	}
}

func TestTypeExcludeIfNotInVersion(t *testing.T) {
	t.Parallel()

	p := Product{
		Name: "test",
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
	}

	cases := []struct {
		description string
		obj         Type
		input       *product.Version
		expected    bool
	}{
		{
			description: "type has Exclude true",
			obj: Type{
				Name:       "test",
				Exclude:    true,
				MinVersion: "",
				ResourceMetadata: &Resource{
					Name:            "test",
					MinVersion:      "",
					ProductMetadata: &p,
				},
			},
			input: &product.Version{
				Name: "ga",
			},
			expected: true,
		},
		{
			description: "type has Exclude false and not empty ExactVersion",
			obj: Type{
				Name:         "test",
				MinVersion:   "",
				Exclude:      false,
				ExactVersion: "beta",
				ResourceMetadata: &Resource{
					Name:            "test",
					MinVersion:      "beta",
					ProductMetadata: &p,
				},
			},
			input: &product.Version{
				Name: "ga",
			},
			expected: true,
		},
		{
			description: "type has Exclude false and empty ExactVersion",
			obj: Type{
				Name:         "test",
				MinVersion:   "beta",
				Exclude:      false,
				ExactVersion: "",
				ResourceMetadata: &Resource{
					Name:            "test",
					MinVersion:      "",
					ProductMetadata: &p,
				},
			},
			input: &product.Version{
				Name: "ga",
			},
			expected: true,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			tc.obj.ExcludeIfNotInVersion(tc.input)
			if got, want := tc.obj.Exclude, tc.expected; !reflect.DeepEqual(got, want) {
				t.Errorf("expected %v to be %v", got, want)
			}
		})
	}
}
