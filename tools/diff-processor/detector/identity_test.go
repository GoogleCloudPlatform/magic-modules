package detector

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDetectMissingIdentityCoverage(t *testing.T) {
	cases := []struct {
		name            string
		resourceContent string
		testContent     string
		resourceName    string
		want            map[string]*MissingIdentityInfo
	}{
		{
			name:            "resource with identity and full coverage is not flagged",
			resourceContent: fakeResourceWithFullCoverage(),
			testContent:     fakeTestWithImportIdentity(),
			resourceName:    "foo_instance",
			want:            map[string]*MissingIdentityInfo{},
		},
		{
			name:            "resource with identity missing Create and Update CRUD and import test is flagged",
			resourceContent: fakeResourceMissingCreateUpdate(),
			testContent:     "",
			resourceName:    "foo_instance",
			want: map[string]*MissingIdentityInfo{
				"foo_instance": {
					MissingCRUD:       []string{"Create", "Update"},
					MissingImportTest: true,
				},
			},
		},
		{
			name:            "resource without identity block is not flagged",
			resourceContent: fakeResourceWithoutIdentity(),
			testContent:     "",
			resourceName:    "foo_instance",
			want:            map[string]*MissingIdentityInfo{},
		},
		{
			name:            "resource file does not exist is skipped",
			resourceContent: "",
			testContent:     "",
			resourceName:    "nonexistent_resource",
			want:            map[string]*MissingIdentityInfo{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			if tc.resourceContent != "" {
				err := os.WriteFile(filepath.Join(tmpDir, "resource_"+tc.resourceName+".go"), []byte(tc.resourceContent), 0644)
				if err != nil {
					t.Fatalf("failed to write resource file: %v", err)
				}
			}
			if tc.testContent != "" {
				err := os.WriteFile(filepath.Join(tmpDir, "resource_"+tc.resourceName+"_test.go"), []byte(tc.testContent), 0644)
				if err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			}

			results, err := DetectMissingIdentityCoverage(tmpDir, []string{tc.resourceName})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(results, tc.want) {
				t.Errorf("got %+v, want %+v", results, tc.want)
			}
		})
	}
}

func TestDetectMissingIdentityCoverage_SkipsUnchangedResources(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "resource_dns_record_set.go"), []byte(fakeResourceMissingCreateUpdate()), 0644)
	if err != nil {
		t.Fatalf("failed to write resource file: %v", err)
	}

	results, err := DetectMissingIdentityCoverage(tmpDir, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected no results when no resources changed, got %d", len(results))
	}
}

// fakeResourceWithFullCoverage returns a resource file that declares ResourceIdentity
// and calls SetResourceIdentityAttributes in Create, Read, and Update.
func fakeResourceWithFullCoverage() string {
	return `package fake

func resourceFooInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	tpgresource.SetResourceIdentityAttributes(d, map[string]interface{}{})
	return nil
}

func resourceFooInstanceRead(d *schema.ResourceData, meta interface{}) error {
	tpgresource.SetResourceIdentityAttributes(d, map[string]interface{}{})
	return nil
}

func resourceFooInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	tpgresource.SetResourceIdentityAttributes(d, map[string]interface{}{})
	return nil
}

func ResourceFooInstance() *schema.Resource {
	return &schema.Resource{
		Identity: &schema.ResourceIdentity{},
	}
}
`
}

// fakeResourceMissingCreateUpdate returns a resource file that declares ResourceIdentity
// but only calls SetResourceIdentityAttributes in Read, leaving Create and Update uncovered.
func fakeResourceMissingCreateUpdate() string {
	return `package fake

func resourceFooInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceFooInstanceRead(d *schema.ResourceData, meta interface{}) error {
	tpgresource.SetResourceIdentityAttributes(d, map[string]interface{}{})
	return nil
}

func resourceFooInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func ResourceFooInstance() *schema.Resource {
	return &schema.Resource{
		Identity: &schema.ResourceIdentity{},
	}
}
`
}

// fakeResourceWithoutIdentity returns a resource file with no ResourceIdentity block.
func fakeResourceWithoutIdentity() string {
	return `package fake

func resourceFooInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceFooInstanceRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func ResourceFooInstance() *schema.Resource {
	return &schema.Resource{}
}
`
}

// fakeTestWithImportIdentity returns a test file that exercises ImportBlockWithResourceIdentity.
func fakeTestWithImportIdentity() string {
	return `package fake

func TestAccFooInstance_importBlockWithResourceIdentity(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ImportStateKind: resource.ImportBlockWithResourceIdentity,
			},
		},
	})
}
`
}
