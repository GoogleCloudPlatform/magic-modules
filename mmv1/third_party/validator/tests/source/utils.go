package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/stretchr/testify/require"
)

func defaultCompareConverterOutput(t *testing.T, expected []google.Asset, actual []google.Asset, offline bool) {
	expectedAssets := normalizeAssets(t, expected, offline)
	actualAssets := normalizeAssets(t, actual, offline)
	require.ElementsMatch(t, expectedAssets, actualAssets)
}

func testConvertCommand(t *testing.T, dir, name string, offline bool, compare compareConvertOutputFunc) {

	if compare == nil {
		compare = defaultCompareConverterOutput
	}

	// Load expected assets
	var expectedRaw []byte
	testfile := filepath.Join(dir, name+".json")
	expectedRaw, err := ioutil.ReadFile(testfile)
	if err != nil {
		t.Fatalf("Error reading %v: %v", testfile, err)
	}
	var expected []google.Asset
	if err := json.Unmarshal(expectedRaw, &expected); err != nil {
		t.Fatalf("unmarshaling: %v", err)
	}

	// Get converted assets
	var actualRaw []byte
	fileNameToConvert := name + ".tfplan.json"
	actualRaw = tfvConvert(t, dir, fileNameToConvert, offline)
	var actual []google.Asset
	err = json.Unmarshal(actualRaw, &actual)
	if err != nil {
		t.Fatalf("unmarshaling: %v", err)
	}

	compare(t, expected, actual, offline)
}

func testValidateCommandGeneric(t *testing.T, dir, name string, offline bool) {

	wantViolation := true
	wantContents := "Constraint GCPAlwaysViolatesConstraintV1.always_violates_all on resource"
	constraintName := "always_violate"

	testValidateCommand(t, wantViolation, wantContents, dir, name, offline, constraintName)
}

func testValidateCommand(t *testing.T, wantViolation bool, wantContents, dir, name string, offline bool, constraintName string) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("cannot get current directory: %v", err)
	}
	policyPath := filepath.Join(cwd, samplePolicyPath, constraintName)
	var got []byte
	got = tfvValidate(t, wantViolation, dir, name+".tfplan.json", policyPath, offline)
	wantRe := regexp.MustCompile(wantContents)
	if wantContents != "" && !wantRe.Match(got) {
		t.Fatalf("binary did not return expect output, \ngot=%s \nwant (regex)=%s", string(got), wantContents)
	}
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
	cmd.Env = []string{"HOME=" + filepath.Join(dir, "fakehome")}
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

func tfvConvert(t *testing.T, dir, tfPlanFile string, offline bool) []byte {
	executable := tfvBinary
	wantError := false
	args := []string{"convert", "--project", data.Provider["project"]}
	if offline {
		args = append(args, "--offline", "--ancestry", data.Ancestry)
	}
	args = append(args, tfPlanFile)
	cmd := exec.Command(executable, args...)
	// Remove environment variables inherited from the test runtime.
	cmd.Env = []string{}
	// Add credentials back.
	if data.Provider["credentials"] != "" {
		cmd.Env = append(cmd.Env, "GOOGLE_APPLICATION_CREDENTIALS="+data.Provider["credentials"])
	}
	cmd.Dir = dir
	payload, _ := run(t, cmd, wantError)
	return payload
}

func tfvValidate(t *testing.T, wantError bool, dir, tfplan, policyPath string, offline bool) []byte {
	executable := tfvBinary
	args := []string{"validate", "--project", data.Provider["project"], "--policy-path", policyPath}
	if offline {
		args = append(args, "--offline", "--ancestry", data.Ancestry)
	}
	args = append(args, tfplan)
	cmd := exec.Command(executable, args...)
	cmd.Env = []string{"GOOGLE_APPLICATION_CREDENTIALS=" + data.Provider["credentials"]}
	cmd.Dir = dir
	payload, _ := run(t, cmd, wantError)
	return payload
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
