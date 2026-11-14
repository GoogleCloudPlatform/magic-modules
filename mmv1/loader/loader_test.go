package loader

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
)

func TestValidateAllProducts(t *testing.T) {
	_, testFilePath, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Failed to get current test file path")
	}
	baseDir := filepath.Dir(filepath.Dir(testFilePath))
	productsDir := filepath.Join(baseDir, "products")
	if _, err := os.Stat(productsDir); os.IsNotExist(err) {
		t.Skipf("products directory not found at %s", productsDir)
	}

	for _, version := range []string{"ga", "beta", "private"} {
		version := version
		t.Run(version, func(t *testing.T) {
			ofs, err := google.NewOverlayFS("", baseDir)
			if err != nil {
				t.Fatalf("NewOverlayFS: %v", err)
			}
			wrappedFS := NewVarsReplacingFS(ofs)
			l := NewLoader(Config{
				Version:        version,
				BaseDirectory:  baseDir,
				Sysfs:          wrappedFS,
				CompilerTarget: "default (terraform)",
			})
			l.LoadProducts()
			if err := l.AddExtraFields(); err != nil {
				t.Fatalf("AddExtraFields: %v", err)
			}
			if errs := l.ValidateErrors(); len(errs) > 0 {
				t.Errorf("ValidateErrors(%s) failed:\n%v", version, errors.Join(errs...))
			}
		})
	}
}
