package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

// testAsset is similar to Asset but with AncestryPath.
type testAsset struct {
	google.Asset
	Ancestry string `json:"ancestry_path"`
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
	credentials := GetTestCredsFromEnv()
	project := GetTestProjectFromEnv()
	org := GetTestOrgFromEnv(nil)
	billingAccount := GetTestBillingAccountFromEnv(nil)
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
			asset.Ancestors = nil
			// remove the parent as the value of that is dependent on project.
			if asset.Resource != nil {
				asset.Resource.Parent = ""
			}
		}
		// Replace placeholder in names. This allows us to compare generated placeholders
		// (for example due to "unknown after apply") with the values in the expected
		// output files.
		asset.Name = re.ReplaceAllString(asset.Name, "/placeholder-foobar")
		if asset.Resource != nil && asset.Resource.Data != nil {
			if _, ok := asset.Resource.Data["projectId"]; ok {
				projectID, _ := asset.Resource.Data["projectId"].(string)
				if strings.HasPrefix(projectID, "placeholder-") {
					asset.Resource.Data["projectId"] = "placeholder-foobar"
				}
			}
			if _, ok := asset.Resource.Data["name"]; ok {
				name, _ := asset.Resource.Data["name"].(string)
				asset.Resource.Data["name"] = re.ReplaceAllString(name, "/placeholder-foobar")
			}
		}
		ret[i] = asset
	}
	return ret
}

func ancestryPathToAncestors(s string) ([]string, error) {
	path := formatAncestryPath(s)
	fragments := strings.Split(path, "/")
	if len(fragments)%2 != 0 {
		return nil, fmt.Errorf("unexpected format of ancestry path: %s", s)
	}
	ancestors := make([]string, len(fragments)/2)
	for i := 0; i < len(ancestors); i++ {
		ancestors[i] = fmt.Sprintf("%s/%s", fragments[i*2], fragments[i*2+1])
	}
	for i, j := 0, len(ancestors)-1; i < j; i, j = i+1, j-1 {
		ancestors[i], ancestors[j] = ancestors[j], ancestors[i]
	}
	return ancestors, nil
}

func formatAncestryPath(s string) string {
	ret := s
	for _, r := range []struct {
		old string
		new string
	}{
		{"organization/", "organizations/"},
		{"folder/", "folders/"},
		{"project/", "projects/"},
	} {
		ret = strings.ReplaceAll(ret, r.old, r.new)
	}
	return ret
}

func readExpectedTestFile(f string) ([]google.Asset, error) {
	payload, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("cannot open %s, got: %s", f, err)
	}
	var want []testAsset
	if err := json.Unmarshal(payload, &want); err != nil {
		return nil, fmt.Errorf("cannot unmarshal JSON into assets: %s", err)
	}
	for ix := range want {
		ancestors, err := ancestryPathToAncestors(want[ix].Ancestry)
		if err != nil {
			return nil, fmt.Errorf("failed to convert to ancestors: %s", err)
		}
		want[ix].Ancestors = ancestors
	}

	ret := make([]google.Asset, len(want))
	for ix := range want {
		ret[ix] = want[ix].Asset
	}
	return ret, nil
}
