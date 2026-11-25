package resourcemanager

import "testing"

func TestIsDefaultServiceAccount(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		email string
		want  bool
	}{
		"compute default SA":                  {"123456789012-compute@developer.gserviceaccount.com", true},
		"app engine default SA":               {"my-project@appspot.gserviceaccount.com", true},
		"app engine SA of domain-scoped ID":   {"my-proj.example.com@appspot.gserviceaccount.com", true},
		"custom SA":                           {"my-sa@my-project.iam.gserviceaccount.com", false},
		"custom SA named like a default":      {"compute-engine-default@my-project.iam.gserviceaccount.com", false},
		"custom SA in another domain":         {"my-project@appspot.example.com", false},
		"appspot as subdomain":                {"my-project@evil.appspot.gserviceaccount.com", false},
		"compute-like address without number": {"x-compute@developer.gserviceaccount.com", false},
		"empty":                               {"", false},
	}

	for name, tc := range cases {
		if got := isDefaultServiceAccount(tc.email); got != tc.want {
			t.Errorf("%s: isDefaultServiceAccount(%q) = %t, want %t", name, tc.email, got, tc.want)
		}
	}
}
