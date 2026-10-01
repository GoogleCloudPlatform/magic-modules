package provider

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/resource"
)

func TestEphemeralArtifactGeneration(t *testing.T) {
	for _, version := range []string{"ga", "beta"} {
		t.Run(version, func(t *testing.T) {
			r, _ := newEphemeralTestResource()
			r.ProductMetadata.Name = "testproduct"
			r.ProductMetadata.ApiName = "testproduct"
			r.ProductMetadata.DisplayName = "Test Product"
			r.ProductMetadata.Versions = []*product.Version{{Name: "ga"}, {Name: "beta"}}
			r.TargetVersionName = version
			r.SourceYamlFile = "products/testproduct/TestResource.yaml"
			r.References = resource.ReferenceLinks{}
			fixture := func(name, minVersion string) *resource.Sample {
				return &resource.Sample{
					Name: name, MinVersion: minVersion, PrimaryResourceId: "test",
					Steps: []*resource.Step{{Name: "read", ConfigPath: "read.tf.tmpl", EphemeralAssertions: map[string]string{
						"has_value": "length(ephemeral.google_test_product_test_resource.test.name) > 0",
					}}},
				}
			}
			r.Ephemeral = &resource.Ephemeral{Generate: true, Samples: []*resource.Sample{fixture("basic", "ga"), fixture("beta", "beta")}}
			r.TestSampleSetUp(fstest.MapFS{"read.tf.tmpl": {Data: []byte(`ephemeral "google_test_product_test_resource" "test" { name = "existing" }`)}})
			if len(r.Samples) != 0 {
				t.Fatal("ephemeral samples leaked into managed-resource samples")
			}
			generator := NewTerraform(r.ProductMetadata, version, time.Now(), os.DirFS(".."))
			r.ImportPath = ImportPathFromVersion(version)
			for _, mode := range []struct {
				name       string
				code, docs bool
			}{
				{"code", true, false}, {"docs", false, true}, {"both", true, true},
			} {
				t.Run(mode.name, func(t *testing.T) {
					out := t.TempDir()
					td := NewTemplateData(out, version, os.DirFS(".."))
					generator.GenerateEphemeralResource(*r, *td, out, mode.code, mode.docs)
					testPath := filepath.Join(out, generator.FolderName(), "services/testproduct", "ephemeral_"+generator.ResourceGoFilename(*r)+"_generated_test.go")
					docPath := filepath.Join(out, "website/docs/ephemeral-resources", generator.FullResourceName(*r)+".html.markdown")
					for path, expected := range map[string]bool{testPath: mode.code, docPath: mode.docs} {
						b, err := os.ReadFile(path)
						if !expected {
							if !os.IsNotExist(err) {
								t.Fatalf("unexpected output: %s", path)
							}
							continue
						}
						if err != nil {
							t.Fatal(err)
						}
						text := string(b)
						if path == testPath {
							if _, err := parser.ParseFile(token.NewFileSet(), path, b, parser.AllErrors); err != nil {
								t.Fatal(err)
							}
							if strings.Contains(text, "_betaGenerated") != (version == "beta") {
								t.Fatal("incorrect sample version filtering")
							}
							for _, want := range []string{"_basicGenerated", `"data.has_value", "true"`, r.ImportPath + "/acctest"} {
								if !strings.Contains(text, want) {
									t.Fatalf("missing %q in test", want)
								}
							}
						} else {
							arguments := strings.Split(strings.Split(text, "## Argument Reference")[1], "## Attributes Reference")[0]
							if !strings.Contains(arguments, "`name` - (Required)") || strings.Contains(arguments, "edition_config") {
								t.Fatal("docs use managed-resource input roles")
							}
							if !strings.Contains(text, "effective_intelligence_config.inner.inner_value") {
								t.Fatal("missing nested output documentation")
							}
							if strings.Contains(text, "### beta") != (version == "beta") {
								t.Fatal("incorrect doc example filtering")
							}
						}
					}
				})
			}
			r.Ephemeral.ExcludeTest = true
			outExcluded := t.TempDir()
			generator.GenerateEphemeralResource(*r, *NewTemplateData(outExcluded, version, os.DirFS("..")), outExcluded, true, true)
			matches, err := filepath.Glob(filepath.Join(outExcluded, generator.FolderName(), "services/testproduct/*_test.go"))
			if err != nil || len(matches) != 0 {
				t.Fatal("exclude_test emitted acceptance tests")
			}
			r.Ephemeral.Generate = false
			out := t.TempDir()
			generator.GenerateEphemeralResource(*r, *NewTemplateData(out, version, os.DirFS("..")), out, true, true)
			entries, err := os.ReadDir(out)
			if err != nil || len(entries) != 0 {
				t.Fatal("disabled generation emitted artifacts")
			}
		})
	}
}

func TestEphemeralInputIdentifierNormalization(t *testing.T) {
	r := &api.Resource{SelfLink: "objects/{{key_id}}"}
	p := &api.Type{Name: "keyId", ResourceMetadata: r}
	if !p.IsEphemeralInput() {
		t.Fatal("snake_case URL identifier must match camelCase property")
	}
}
