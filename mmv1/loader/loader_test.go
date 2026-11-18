package loader_test

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/loader"
)

func TestValidateAllProducts(t *testing.T) {
	_, testFilePath, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve test file path")
	}
	mmv1Dir := filepath.Dir(filepath.Dir(testFilePath))
	if _, err := os.Stat(filepath.Join(mmv1Dir, "products")); os.IsNotExist(err) {
		t.Skipf("products directory not found at %s (sandboxed test run)", mmv1Dir)
	}

	ofs, err := google.NewOverlayFS("", mmv1Dir)
	if err != nil {
		t.Fatalf("NewOverlayFS failed: %v", err)
	}
	wrappedFS := loader.NewVarsReplacingFS(ofs)

	l := loader.NewLoader(loader.Config{
		Version:        "internal",
		BaseDirectory:  mmv1Dir,
		Sysfs:          wrappedFS,
		CompilerTarget: "terraform",
	})
	l.LoadProducts()
	if err := l.AddExtraFields(); err != nil {
		t.Fatalf("AddExtraFields failed: %v", err)
	}

	productKeys := make([]string, 0, len(l.Products))
	for k := range l.Products {
		productKeys = append(productKeys, k)
	}
	sort.Strings(productKeys)

	for _, pk := range productKeys {
		product := l.Products[pk]
		ver := product.VersionObjOrClosest("internal")
		for _, res := range product.Objects {
			res.ExcludeIfNotInVersion(ver)
			for _, err := range res.Validate() {
				t.Errorf("%s: %v", res.SourceYamlFile, err)
			}
		}
	}
}
