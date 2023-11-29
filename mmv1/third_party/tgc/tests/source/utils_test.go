package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai"
	"github.com/google/go-cmp/cmp"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
	"go.uber.org/zap/zaptest"
)

func defaultCompareConverterOutput(t *testing.T, expected []caiasset.Asset, actual []caiasset.Asset, offline bool) {
	expectedAssets := normalizeAssets(t, expected, offline)
	actualAssets := normalizeAssets(t, actual, offline)
	if diff := cmp.Diff(expectedAssets, actualAssets); diff != "" {
		t.Errorf("%v diff(-want, +got):\n%s", t.Name(), diff)
	}
}

func testConvertCommand(t *testing.T, dir, tfplanName string, jsonName string, offline bool, withProject bool, compare compareConvertOutputFunc) {

	if compare == nil {
		compare = defaultCompareConverterOutput
	}

	// Load expected assets
	expected, err := readExpectedTestFile(filepath.Join(dir, jsonName+".json"))
	if err != nil {
		t.Fatal(err)
	}

	// Get converted assets
	fileNameToConvert := tfplanName + ".tfplan.json"
	actual := tfvConvert(t, dir, fileNameToConvert, offline, withProject)

	compare(t, expected, actual, offline)
}

func terraformWorkflow(t *testing.T, dir, name string) {
	terraformInit(t, "terraform", dir)
	terraformPlan(t, "terraform", dir, name+".tfplan")
	payload := terraformShow(t, "terraform", dir, name+".tfplan")
	saveFile(t, dir, name+".tfplan.json", payload)
}

func terraformInit(t *testing.T, executable, dir string) {
	terraformExec(t, executable, dir, "init", "-input=false")
}

func terraformPlan(t *testing.T, executable, dir, tfplan string) {
	terraformExec(t, executable, dir, "plan", "-input=false", "-refresh=false", "-out", tfplan)
}

func terraformShow(t *testing.T, executable, dir, tfplan string) []byte {
	return terraformExec(t, executable, dir, "show", "--json", tfplan)
}

func terraformExec(t *testing.T, executable, dir string, args ...string) []byte {
	cmd := exec.Command(executable, args...)
	cmd.Env = []string{
		"HOME=" + filepath.Join(dir, "fakehome"),
		"GOOGLE_PROJECT=" + data.Provider["project"],
		"GOOGLE_FOLDER=" + data.FolderID,
		"GOOGLE_ORG=" + data.OrgID,
		"GOOGLE_OAUTH_ACCESS_TOKEN=fake-token", // GOOGLE_OAUTH_ACCESS_TOKEN is required so terraform plan does not require the google authentication cert
	}
	if os.Getenv("TF_CLI_CONFIG_FILE") != "" {
		cmd.Env = append(cmd.Env, "TF_CLI_CONFIG_FILE="+os.Getenv("TF_CLI_CONFIG_FILE"))
	}
	cmd.Dir = dir
	wantError := false
	payload, _ := run(t, cmd, wantError)
	return payload
}

func saveFile(t *testing.T, dir, filename string, payload []byte) {
	fullpath := filepath.Join(dir, filename)
	f, err := os.Create(fullpath)
	if err != nil {
		t.Fatalf("error while creating file %s, error %v", fullpath, err)
	}
	_, err = f.Write(payload)
	if err != nil {
		t.Fatalf("error while writing to file %s, error %v", fullpath, err)
	}
}

func tfvConvert(t *testing.T, dir, tfPlanFile string, offline bool, withProject bool) []caiasset.Asset {
	jsonPlan, err := os.ReadFile(filepath.Join(dir, tfPlanFile))
	if err != nil {
		t.Fatalf("Error parsing %s: %s", tfPlanFile, err)
	}
	ctx := context.Background()
	opts := &tfplan2cai.Options{
		ConvertUnchanged: false,
		ErrorLogger:      zaptest.NewLogger(t),
		Offline:          offline,
		DefaultRegion:    "",
		DefaultZone:      "",
		UserAgent:        "",
	}
	if withProject {
		opts.DefaultProject = data.Provider["project"]
		opts.AncestryCache = map[string]string{
			data.Provider["project"]: data.Ancestry,
		}
	}
	got, err := tfplan2cai.Convert(ctx, jsonPlan, opts)
	if err != nil {
		t.Fatalf("Error converting %s: %s", tfPlanFile, err)
	}
	return got
}

// run a command and call t.Fatal on non-zero exit.
func run(t *testing.T, cmd *exec.Cmd, wantError bool) ([]byte, []byte) {
	var stderr, stdout bytes.Buffer
	cmd.Stderr, cmd.Stdout = &stderr, &stdout
	err := cmd.Run()
	if gotError := (err != nil); gotError != wantError {
		t.Fatalf("running %s: \nerror=%v \nstderr=%s \nstdout=%s", cmdToString(cmd), err, stderr.String(), stdout.String())
	}
	// Print env, stdout and stderr if verbose flag is used.
	if len(cmd.Env) != 0 {
		t.Logf("=== Environment Variable of %s ===", cmdToString(cmd))
		t.Log(strings.Join(cmd.Env, "\n"))
	}
	if stdout.String() != "" {
		t.Logf("=== STDOUT of %s ===", cmdToString(cmd))
		t.Log(stdout.String())
	}
	if stderr.String() != "" {
		t.Logf("=== STDERR of %s ===", cmdToString(cmd))
		t.Log(stderr.String())
	}
	return stdout.Bytes(), stderr.Bytes()
}

// cmdToString clones the logic of https://golang.org/pkg/os/exec/#Cmd.String.
func cmdToString(c *exec.Cmd) string {
	// report the exact executable path (plus args)
	b := new(strings.Builder)
	b.WriteString(c.Path)
	for _, a := range c.Args[1:] {
		b.WriteByte(' ')
		b.WriteString(a)
	}
	return b.String()
}

func generateTFVconvertedAsset(t *testing.T, testDir, testSlug string) {
	// Get converted assets
	fileNameToConvert := testSlug + ".tfplan.json"
	assets := tfvConvert(t, testDir, fileNameToConvert, true, true)
	dstDir := "../testdata/generatedconvert"
	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		os.MkdirAll(dstDir, 0700)
	}

	conversionRaw, err := json.Marshal(assets)
	if err != nil {
		t.Fatalf("Failed to convert asset: %s", err)
	}

	conversionPretty := &bytes.Buffer{}
	if err := json.Indent(conversionPretty, conversionRaw, "", "  "); err != nil {
		panic(err)
	}

	dstFile := path.Join(dstDir, testSlug+".json")
	err = os.WriteFile(dstFile, conversionPretty.Bytes(), 0666)
	if err != nil {
		t.Fatalf("error while writing to file %s, error %v", dstFile, err)
	}

	fmt.Println("created file : " + dstFile)
}

// newTestConfig create a config using the http test server.
func newTestConfig(server *httptest.Server) *transport_tpg.Config {
	cfg := &transport_tpg.Config{}
	cfg.Client = server.Client()
	configureTestBasePaths(cfg, server.URL)
	return cfg
}

func configureTestBasePaths(c *transport_tpg.Config, url string) {
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	typ := reflect.ValueOf(c).Elem().Type()
	val := reflect.ValueOf(c).Elem()

	for i := 0; i < typ.NumField(); i++ {
		name := typ.Field(i).Name
		if strings.HasSuffix(name, "BasePath") {
			val.Field(i).SetString(url)
		}
	}
}
