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
				BaseUrl: "https://abc.googleapis.com/beta/",
			},
			expected: "abc.googleapis.com",
		},
		{
			description: "BaseUrl with locational subdomain",
			obj: Product{
				BaseUrl: "https://{{location}}-abc.googleapis.com/ga/",
			},
			expected: "abc.googleapis.com",
		},
		{
			description: "BaseUrl and CaiBaseUrl",
			obj: Product{
				BaseUrl:    "https://abc.googleapis.com/ga/",
				CaiBaseUrl: "https://def.googleapis.com/ga/",
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
				BaseUrl: "https://abc.googleapis.com/v1/",
			},
			expected: "v1",
		},
		{
			description: "BaseUrl without trailing /",
			obj: Product{
				BaseUrl: "https://abc.googleapis.com/v1",
			},
			expected: "v1",
		},
		{
			description: "BaseUrl with version of beta",
			obj: Product{
				BaseUrl: "https://abc.googleapis.com/beta/",
			},
			expected: "beta",
		},
		{
			description: "BaseUrl without valid version",
			obj: Product{
				BaseUrl: "https://abc.googleapis.com/other/",
			},
			expected: "",
		},
		{
			description: "BaseUrl with additional value in path",
			obj: Product{
				BaseUrl: "https://abc.googleapis.com/compute/v1/",
			},
			expected: "v1",
		},
		{
			description: "standard BaseUrl",
			obj: Product{
				BaseUrl:    "https://{{location}}-abc.googleapis.com/",
				CaiBaseUrl: "https://abc.googleapis.com/v1/",
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
