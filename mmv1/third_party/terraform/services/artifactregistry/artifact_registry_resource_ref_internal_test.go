package artifactregistry

import "testing"

func TestArtifactResourceRef(t *testing.T) {
	cases := []struct {
		name    string
		refName string
		sep     string
		version string
		want    string
	}{
		{
			name:    "npm semver left intact",
			refName: "my-package",
			sep:     ":",
			version: "1.2.3",
			want:    "my-package:1.2.3",
		},
		{
			name:    "docker digest keeps its colon",
			refName: "my-image",
			sep:     "@",
			version: "sha256:2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
			want:    "my-image@sha256:2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		{
			name:    "no version escapes only the name",
			refName: "krane/debug",
			sep:     ":",
			version: "",
			want:    "krane%2Fdebug",
		},
		{
			name:    "traversal in version is neutralized",
			refName: "my-package",
			sep:     ":",
			version: "1.0.0/../../otherRepo/npmPackages/evil",
			want:    "my-package:1.0.0%2F..%2F..%2FotherRepo%2FnpmPackages%2Fevil",
		},
		{
			name:    "query and fragment in version are neutralized",
			refName: "my-image",
			sep:     "@",
			version: "sha256:abc?alt=media#frag",
			want:    "my-image@sha256:abc%3Falt=media%23frag",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := artifactResourceRef(tc.refName, tc.sep, tc.version)
			if got != tc.want {
				t.Errorf("artifactResourceRef(%q, %q, %q) = %q, want %q", tc.refName, tc.sep, tc.version, got, tc.want)
			}
		})
	}
}
