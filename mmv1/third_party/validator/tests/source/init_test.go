package test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
)

var (
	data      *testData
	tfvBinary string
	tmpDir    = os.TempDir()
)

// testData represents the full dataset that is used for templating terraform
// configs. It contains Google API resources that are expected to be returned
// after converting the terraform plan.
type testData struct {
	// is not nil - Terraform 12 version used
	TFVersion string
	// provider "google"
	Provider map[string]string
	Project  map[string]string
	Time     map[string]string
	OrgID    string
	FolderID string
	Ancestry string
}

// init initializes the variables used for testing. As tests rely on
// environment variables, the parsing of those are only done once.
func init() {
	// don't raise errors in glog

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot get current directory: %v", err)
	}
	tfvBinary = filepath.Join(cwd, "..", "bin", "terraform-validator")
	credentials := getTestCredsFromEnv()
	project := getTestProjectFromEnv()
	org := getTestOrgFromEnv(nil)
	billingAccount := getTestBillingAccountFromEnv(nil)
	folder, ok := os.LookupEnv("TEST_FOLDER_ID")
	if !ok {
		log.Printf("Missing required env var TEST_FOLDER_ID. Default (%s) will be used.", defaultFolder)
		folder = defaultFolder
	}
	ancestry, ok := os.LookupEnv("TEST_ANCESTRY")
	if !ok {
		log.Printf("Missing required env var TEST_ANCESTRY. Default (%s) will be used.", defaultAncestry)
		ancestry = defaultAncestry
	}
	providerVersion := defaultProviderVersion
	//As time is not information in terraform resource data, time is fixed for testing purposes
	fixedTime := time.Date(2021, time.April, 14, 15, 16, 17, 0, time.UTC)
	data = &testData{
		TFVersion: "0.12",
		Provider: map[string]string{
			"version":     providerVersion,
			"project":     project,
			"credentials": credentials,
		},
		Time: map[string]string{
			"RFC3339Nano": fixedTime.Format(time.RFC3339Nano),
		},
		Project: map[string]string{
			"Name":               "My Project Name",
			"ProjectId":          "my-project-id",
			"BillingAccountName": billingAccount,
			"Number":             "1234567890",
		},
		OrgID:    org,
		FolderID: folder,
		Ancestry: ancestry,
	}
}

func generateTestFiles(t *testing.T, sourceDir string, targetDir string, selector string) {
	funcMap := template.FuncMap{
		"pastLastSlash": func(s string) string {
			split := strings.Split(s, "/")
			return split[len(split)-1]
		},
	}
	tmpls, err := template.New("").Funcs(funcMap).
		ParseGlob(filepath.Join(sourceDir, selector))
	if err != nil {
		t.Fatalf("generateTestFiles: %v", err)
	}
	for _, tmpl := range tmpls.Templates() {
		if tmpl.Name() == "" {
			continue // Skip base template.
		}
		path := filepath.Join(targetDir, tmpl.Name())
		f, err := os.Create(path)
		if err != nil {
			t.Fatalf("creating terraform file %v: %v", path, err)
		}
		if err := tmpl.Execute(f, data); err != nil {
			t.Fatalf("templating terraform file %v: %v", path, err)
		}
		if err := f.Close(); err != nil {
			t.Fatalf("closing file %v: %v", path, err)
		}
		t.Logf("Successfully created file %v", path)
	}
}

func normalizeAssets(t *testing.T, assets []google.Asset, offline bool) []google.Asset {
	t.Helper()
	ret := make([]google.Asset, len(assets))
	re := regexp.MustCompile(`/placeholder-[^/]+`)
	for i := range assets {
		// Get conformity by converting to/from json.
		bytes, err := json.Marshal(assets[i])
		if err != nil {
			t.Fatalf("marshaling: %v", err)
		}

		var asset google.Asset
		err = json.Unmarshal(bytes, &asset)
		if err != nil {
			t.Fatalf("marshaling: %v", err)
		}
		if !offline {
			// remove the ancestry as the value of that is dependent on project,
			// and is not important for the test.
			asset.Ancestry = ""
		}
		// Replace placeholder in names. This allows us to compare generated placeholders
		// (for example due to "unknown after apply") with the values in the expected
		// output files.
		asset.Name = re.ReplaceAllString(asset.Name, fmt.Sprintf("/placeholder-foobar"))
		ret[i] = asset
	}
	return ret
}
