package resourcemanager

import "testing"

func TestIsDefaultServiceAccount(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		email     string
		projectID string
		want      bool
	}{
		"compute default SA":                  {"123456789012-compute@developer.gserviceaccount.com", "my-project", true},
		"app engine default SA":               {"my-project@appspot.gserviceaccount.com", "my-project", true},
		"app engine SA of another project":    {"other-project@appspot.gserviceaccount.com", "my-project", false},
		"custom SA":                           {"my-sa@my-project.iam.gserviceaccount.com", "my-project", false},
		"custom SA named like a default":      {"compute-engine-default@my-project.iam.gserviceaccount.com", "my-project", false},
		"compute-like address without number": {"x-compute@developer.gserviceaccount.com", "my-project", false},
		"empty":                               {"", "my-project", false},
	}

	for name, tc := range cases {
		if got := isDefaultServiceAccount(tc.email, tc.projectID); got != tc.want {
			t.Errorf("%s: isDefaultServiceAccount(%q, %q) = %t, want %t", name, tc.email, tc.projectID, got, tc.want)
		}
	}
}
