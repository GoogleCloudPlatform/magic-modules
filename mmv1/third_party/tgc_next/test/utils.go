package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func terraformWorkflow(dir, name, project string) error {
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
			return fmt.Errorf("Error writing provider override to %s: %v", tfFile, err)
		}
	}

	if err := terraformPlan("terraform", dir, project, name+".tfplan"); err != nil {
		return err
	}
	payload, err := terraformShow("terraform", dir, project, name+".tfplan")
	if err != nil {
		return err
	}
	if err := saveFile(dir, name+".tfplan.json", payload); err != nil {
		return err
	}
	return nil
}

func terraformInit(executable, dir, project string) error {
	_, err := terraformExec(executable, dir, project, "init", "-input=false")
	return err
}

func terraformPlan(executable, dir, project, tfplan string) error {
	_, err := terraformExec(executable, dir, project, "plan", "-input=false", "-refresh=false", "-out", tfplan)
	return err
}

func terraformShow(executable, dir, project, tfplan string) ([]byte, error) {
	return terraformExec(executable, dir, project, "show", "--json", tfplan)
}

func terraformExec(executable, dir, project string, args ...string) ([]byte, error) {
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
	payload, _, err := run(cmd, wantError)
	return payload, err
}

func saveFile(dir, filename string, payload []byte) error {
	fullpath := filepath.Join(dir, filename)
	f, err := os.Create(fullpath)
	if err != nil {
		return fmt.Errorf("error while creating file %s, error %v", fullpath, err)
	}
	defer f.Close()
	_, err = f.Write(payload)
	if err != nil {
		return fmt.Errorf("error while writing to file %s, error %v", fullpath, err)
	}
	return nil
}

// run a command and return error on non-zero exit instead of t.Fatalf.
func run(cmd *exec.Cmd, wantError bool) ([]byte, []byte, error) {
	var stderr, stdout bytes.Buffer
	cmd.Stderr, cmd.Stdout = &stderr, &stdout
	err := cmd.Run()

	// Do not log output here to avoid cluttering test results on retries.
	// The full output is included in the error returned on failure.

	if gotError := (err != nil); gotError != wantError {
		return stdout.Bytes(), stderr.Bytes(), fmt.Errorf("running %s: \nerror=%v \nstderr=%s", cmd.String(), err, stderr.String())
	}
	return stdout.Bytes(), stderr.Bytes(), nil
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
	_, after, found := strings.Cut(fullTestName, "/")
	if !found {
		return ""
	}
	return after
}
