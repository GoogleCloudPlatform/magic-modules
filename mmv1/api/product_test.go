package api

import (
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
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
