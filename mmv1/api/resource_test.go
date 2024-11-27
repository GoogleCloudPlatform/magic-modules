package api

import (
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
)

func TestResourceMinVersionObj(t *testing.T) {
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
		obj         Resource
		expected    string
	}{
		{
			description: "resource minVersion is empty",
			obj: Resource{
				Name:            "test",
				MinVersion:      "",
				ProductMetadata: &p,
			},
			expected: "ga",
		},
		{
			description: "resource minVersion is not empty",
			obj: Resource{
				Name:            "test",
				MinVersion:      "beta",
				ProductMetadata: &p,
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

func TestResourceNotInVersion(t *testing.T) {
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
		obj         Resource
		input       *product.Version
		expected    bool
	}{
		{
			description: "ga is in version if MinVersion is empty",
			obj: Resource{
				Name:            "test",
				MinVersion:      "",
				ProductMetadata: &p,
			},
			input: &product.Version{
				Name: "ga",
			},
			expected: false,
		},
		{
			description: "ga is not in version if MinVersion is beta",
			obj: Resource{
				Name:            "test",
				MinVersion:      "beta",
				ProductMetadata: &p,
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

			if got, want := tc.obj.NotInVersion(tc.input), tc.expected; !reflect.DeepEqual(got, want) {
				t.Errorf("expected %v to be %v", got, want)
			}
		})
	}
}

func TestResourceServiceVersion(t *testing.T) {
	t.Parallel()

	cases := []struct {
		description string
		obj         Resource
		expected    string
	}{
		{
			description: "BaseUrl does not start with a version",
			obj: Resource{
				BaseUrl: "test",
			},
			expected: "",
		},
		{
			description: "BaseUrl starts with / and does not include a version",
			obj: Resource{
				BaseUrl: "/test",
			},
			expected: "",
		},
		{
			description: "BaseUrl starts with a version",
			obj: Resource{
				BaseUrl: "v3/test",
			},
			expected: "v3",
		},
		{
			description: "BaseUrl starts with a / followed by version",
			obj: Resource{
				BaseUrl: "/v3/test",
			},
			expected: "v3",
		},
		{
			description: "CaiBaseUrl does not start with a version",
			obj: Resource{
				BaseUrl:    "apis/serving.knative.dev/v1/namespaces/{{project}}/services",
				CaiBaseUrl: "projects/{{project}}/locations/{{location}}/services",
			},
			expected: "",
		},
		{
			description: "CaiBaseUrl starts with a version",
			obj: Resource{
				BaseUrl:    "apis/serving.knative.dev/v1/namespaces/{{project}}/services",
				CaiBaseUrl: "v1/projects/{{project}}/locations/{{location}}/services",
			},
			expected: "v1",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			if got, want := tc.obj.ServiceVersion(), tc.expected; got != want {
				t.Errorf("expected %q to be %q", got, want)
			}
		})
	}
}
