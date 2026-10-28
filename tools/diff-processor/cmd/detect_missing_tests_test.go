package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/detector"
	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDetectMissingUpdateTestsCmd(t *testing.T) {
	servicesDir := t.TempDir()
	svcDir := filepath.Join(servicesDir, "example")
	if err := os.MkdirAll(svcDir, 0o755); err != nil {
		t.Fatalf("failed creating service dir: %v", err)
	}
	testSrc := `package example_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
)

func TestAccExample(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: ` + "`" + `
resource "google_example" "primary" {
  unmutated = "same"
  mutated   = "v1"
}
` + "`" + `,
			},
			{
				Config: ` + "`" + `
resource "google_example" "primary" {
  unmutated = "same"
  mutated   = "v2"
}
` + "`" + `,
			},
		},
	})
}
`
	if err := os.WriteFile(filepath.Join(svcDir, "resource_example_test.go"), []byte(testSrc), 0o600); err != nil {
		t.Fatalf("failed writing test file: %v", err)
	}

	var buf bytes.Buffer
	o := &detectMissingTestsOptions{
		computeSchemaDiff: func() diff.SchemaDiff {
			return diff.ComputeSchemaDiff(
				map[string]*schema.Resource{},
				map[string]*schema.Resource{
					"google_example": {
						Schema: map[string]*schema.Schema{
							"unmutated": {Type: schema.TypeString, Optional: true},
							"mutated":   {Type: schema.TypeString, Optional: true},
							"immutable": {Type: schema.TypeString, Optional: true, ForceNew: true},
						},
					},
				},
			)
		},
		stdout: &buf,
	}

	if err := o.runUpdateTests([]string{servicesDir}); err != nil {
		t.Fatalf("runUpdateTests failed: %v", err)
	}

	var got map[string]*detector.MissingUpdateTestInfo
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal json output: %v", err)
	}
	want := map[string]*detector.MissingUpdateTestInfo{
		"google_example": {
			UntestedUpdateFields: []string{"unmutated"},
			Tests:                []string{"TestAccExample"},
		},
	}
	if d := cmp.Diff(want, got); d != "" {
		t.Errorf("unexpected output (-want +got):\n%s", d)
	}
}
