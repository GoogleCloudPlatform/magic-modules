package provider

import (
	"bytes"
	"errors"
	"log"
	"os"
	"path/filepath"
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

func TestCopyDir(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create nested source structure with distinct permissions
	if err := os.MkdirAll(filepath.Join(srcDir, "nested", "sub"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(srcDir, "services", "compute"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "root.txt"), []byte("new-root-content"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "nested", "sub", "script.sh"), []byte("#!/bin/sh\necho hi\n"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "nested", "skip.txt"), []byte("should-skip"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "services", "compute", "svc.go"), []byte("package compute"), 0644); err != nil {
		t.Fatal(err)
	}

	// Pre-populate destination file to verify overwriting
	if err := os.WriteFile(filepath.Join(dstDir, "root.txt"), []byte("old-content"), 0644); err != nil {
		t.Fatal(err)
	}

	err := copyDir(srcDir, dstDir, func(src string) bool {
		return strings.Contains(src, "/services/") || strings.HasSuffix(src, "skip.txt")
	})
	if err != nil {
		t.Fatalf("copyDir failed: %v", err)
	}

	// Verify overwritten root.txt
	gotRoot, err := os.ReadFile(filepath.Join(dstDir, "root.txt"))
	if err != nil {
		t.Fatalf("reading root.txt: %v", err)
	}
	if string(gotRoot) != "new-root-content" {
		t.Errorf("root.txt = %q, want %q", string(gotRoot), "new-root-content")
	}

	// Verify nested file and permissions
	scriptPath := filepath.Join(dstDir, "nested", "sub", "script.sh")
	gotScript, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("reading script.sh: %v", err)
	}
	if string(gotScript) != "#!/bin/sh\necho hi\n" {
		t.Errorf("script.sh = %q, want %q", string(gotScript), "#!/bin/sh\necho hi\n")
	}
	info, err := os.Stat(scriptPath)
	if err != nil {
		t.Fatalf("stat script.sh: %v", err)
	}
	if info.Mode().Perm() != 0755 {
		t.Errorf("script.sh mode = %o, want 0755", info.Mode().Perm())
	}

	// Verify skipped file and skipped directory
	if _, err := os.Stat(filepath.Join(dstDir, "nested", "skip.txt")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected skip.txt to be skipped, stat err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(dstDir, "services", "compute", "svc.go")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected services/compute/svc.go to be skipped, stat err = %v", err)
	}

	// Verify non-existent source returns os.ErrNotExist
	if err := copyDir(filepath.Join(srcDir, "does-not-exist"), dstDir, nil); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected os.ErrNotExist for missing srcDir, got %v", err)
	}
}
