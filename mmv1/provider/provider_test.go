package provider

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
)

func TestTerraformVerboseLogging(t *testing.T) {
	res := &api.Resource{Name: "Instance"}
	p := &api.Product{
		Name:        "Compute",
		Versions:    []*product.Version{{Name: "ga"}},
		Objects:     []*api.Resource{res},
		PackagePath: "products/compute",
	}
	p.Version = p.Versions[0]
	res.SetDefault(p)

	testCases := []struct {
		name         string
		verbose      bool
		expectLogged bool
	}{
		{
			name:         "Default non-verbose suppresses Generating resource log",
			verbose:      false,
			expectLogged: false,
		},
		{
			name:         "Verbose flag enables Generating resource log",
			verbose:      true,
			expectLogged: true,
		},
	}

	fsys := os.DirFS("..")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			google.VerboseLogging = tc.verbose
			t.Cleanup(func() {
				google.VerboseLogging = false
			})

			var buf bytes.Buffer
			log.SetOutput(&buf)
			t.Cleanup(func() {
				log.SetOutput(log.Writer())
			})

			providerObj := NewTerraform(p, "ga", time.Now(), fsys)
			// GenerateObjects will call GenerateObject which logs "Generating Instance resource" if google.VerboseLogging is true
			providerObj.GenerateObjects(t.TempDir(), "", false, false)

			output := buf.String()
			containsLog := strings.Contains(output, "Generating Instance resource")
			if containsLog != tc.expectLogged {
				t.Errorf("expected logged=%v, got output: %q", tc.expectLogged, output)
			}
		})
	}
}

func TestExpectedOutputFolder(t *testing.T) {
	testCases := []struct {
		path     string
		expected bool
	}{
		{"/Users/user/git/terraform-provider-google", true},
		{"/Users/user/git/terraform-provider-google-beta", true},
		{"/usr/local/google/home/user/.gemini/jetski/worktrees/terraform-provider-google-beta/my_branch", true},
		{"/some/random/dir", false},
	}

	for _, tc := range testCases {
		got := expectedOutputFolder(tc.path)
		if got != tc.expected {
			t.Errorf("expectedOutputFolder(%q) = %v; want %v", tc.path, got, tc.expected)
		}
	}
}

func TestIsHashicorpTarget(t *testing.T) {
	testCases := []struct {
		path     string
		expected bool
	}{
		{"/Users/user/git/terraform-provider-google", true},
		{"/Users/user/git/terraform-provider-google-beta", true},
		{"/usr/local/google/home/user/.gemini/jetski/worktrees/terraform-provider-google-beta/my_branch", true},
		{"/Users/user/git/terraform-google-conversion", false},
		{"/usr/local/google/home/user/.gemini/jetski/worktrees/terraform-google-conversion/my_branch", false},
		{"/some/random/dir", false},
	}

	for _, tc := range testCases {
		got := isHashicorpTarget(tc.path)
		if got != tc.expected {
			t.Errorf("isHashicorpTarget(%q) = %v; want %v", tc.path, got, tc.expected)
		}
	}
}
