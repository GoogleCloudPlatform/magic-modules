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
		"GOOGLE_PROJECT=" + defaultProject,
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
