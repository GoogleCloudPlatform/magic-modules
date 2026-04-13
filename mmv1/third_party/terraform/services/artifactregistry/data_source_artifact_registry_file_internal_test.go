package artifactregistry

import "testing"

func TestBuildFileResourceURL(t *testing.T) {
	cases := []struct {
		name     string
		base     string
		project  string
		location string
		repo     string
		fileID   string
		want     string
	}{
		{
			name:     "simple generic file",
			base:     "https://artifactregistry.googleapis.com/v1/",
			project:  "my-proj",
			location: "us-central1",
			repo:     "my-repo",
			fileID:   "foo.tar.gz",
			want:     "https://artifactregistry.googleapis.com/v1/projects/my-proj/locations/us-central1/repositories/my-repo/files/foo.tar.gz",
		},
		{
			name:     "maven file with slashes and colons",
			base:     "https://artifactregistry.googleapis.com/v1/",
			project:  "p",
			location: "us",
			repo:     "r",
			fileID:   "com.google.guava:guava:32.0.0:guava-32.0.0.jar",
			want:     "https://artifactregistry.googleapis.com/v1/projects/p/locations/us/repositories/r/files/com.google.guava%3Aguava%3A32.0.0%3Aguava-32.0.0.jar",
		},
		{
			name:     "path with slashes",
			base:     "https://artifactregistry.googleapis.com/v1/",
			project:  "p",
			location: "us",
			repo:     "r",
			fileID:   "nested/path/file.txt",
			want:     "https://artifactregistry.googleapis.com/v1/projects/p/locations/us/repositories/r/files/nested%2Fpath%2Ffile.txt",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildFileResourceURL(tc.base, tc.project, tc.location, tc.repo, tc.fileID)
			if got != tc.want {
				t.Errorf("buildFileResourceURL() = %q, want %q", got, tc.want)
			}
		})
	}
}
