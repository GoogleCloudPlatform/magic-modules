package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Writes the data into a JSON file
func writeJSONFile(filename string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("Error marshaling data for %s: %v\n", filename, err)
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("Error writing to file %s: %v\n", filename, err)
	}
	return nil
}

const (
	defaultOrganization = "529579013760"
	defaultProject      = "ci-test-project-nightly-beta"
)

func terraformWorkflow(t *testing.T, dir, name, project string) {
	defer os.Remove(filepath.Join(dir, fmt.Sprintf("%s.tf", name)))
	defer os.Remove(filepath.Join(dir, fmt.Sprintf("%s.tfplan", name)))

	// TODO: remove this when we have a proper google provider
	// Inject google-beta provider override
	tfFile := filepath.Join(dir, fmt.Sprintf("%s.tf", name))
	content, err := os.ReadFile(tfFile)
	if err == nil {
		override := `
terraform {
  required_providers {
    google = {
      source = "hashicorp/google-beta"
    }
  }
}
`
		if err := os.WriteFile(tfFile, append(content, []byte(override)...), 0644); err != nil {
			t.Fatalf("Error writing provider override to %s: %v", tfFile, err)
		}
	}

	terraformInit(t, "terraform", dir, project)
	terraformPlan(t, "terraform", dir, project, name+".tfplan")
	payload := terraformShow(t, "terraform", dir, project, name+".tfplan")
	saveFile(t, dir, name+".tfplan.json", payload)
}

func terraformInit(t *testing.T, executable, dir, project string) {
	terraformExec(t, executable, dir, project, "init", "-input=false")
}

func terraformPlan(t *testing.T, executable, dir, project, tfplan string) {
	terraformExec(t, executable, dir, project, "plan", "-input=false", "-refresh=false", "-out", tfplan)
}

func terraformShow(t *testing.T, executable, dir, project, tfplan string) []byte {
	return terraformExec(t, executable, dir, project, "show", "--json", tfplan)
}

func terraformExec(t *testing.T, executable, dir, project string, args ...string) []byte {
	if project == "" {
		project = defaultProject
	}
	cmd := exec.Command(executable, args...)
	cmd.Env = []string{
		"HOME=" + filepath.Join(dir, "fakehome"),
		"GOOGLE_PROJECT=" + project,
		"GOOGLE_FOLDER=" + "",
		"GOOGLE_ORG=" + defaultOrganization,
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

// run a command and call t.Fatal on non-zero exit.
func run(t *testing.T, cmd *exec.Cmd, wantError bool) ([]byte, []byte) {
	var stderr, stdout bytes.Buffer
	cmd.Stderr, cmd.Stdout = &stderr, &stdout
	err := cmd.Run()
	if gotError := (err != nil); gotError != wantError {
		t.Fatalf("running %s: \nerror=%v \nstderr=%s \nstdout=%s", cmd.String(), err, stderr.String(), stdout.String())
	}
	// Print env, stdout and stderr if verbose flag is used.
	if len(cmd.Env) != 0 {
		t.Logf("=== Environment Variable of %s ===", cmd.String())
		t.Log(strings.Join(cmd.Env, "\n"))
	}
	if stdout.String() != "" {
		t.Logf("=== STDOUT of %s ===", cmd.String())
		t.Log(stdout.String())
	}
	if stderr.String() != "" {
		t.Logf("=== STDERR of %s ===", cmd.String())
		t.Log(stderr.String())
	}
	return stdout.Bytes(), stderr.Bytes()
}

// Creates a deep copy of a source map using JSON marshalling and unmarshalling.
func DeepCopyMap(source interface{}, destination interface{}) error {
	marshalled, err := json.Marshal(source)
	if err != nil {
		return fmt.Errorf("failed to marshal source map: %w", err)
	}

	err = json.Unmarshal(marshalled, destination)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON into destination map: %w", err)
	}

	return nil
}

type TestCase struct {
	Name string
	Skip string
}

func GetSubTestName(fullTestName string) string {
	parts := strings.Split(fullTestName, "/")

	// Get the index of the last element
	lastIndex := len(parts) - 1

	// Check for an empty or malformed string
	if lastIndex < 0 {
		return ""
	}

	return parts[lastIndex]
}
