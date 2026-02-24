// Copyright 2025 Google Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resource

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
	"github.com/golang/glog"
)

type Step struct {
	Name string `yaml:"name,omitempty"`

	// The path to this step's Terraform config.
	// Defaults to `templates/terraform/samples/{{service}}/{{name}}.tf.erb`
	ConfigPath string `yaml:"config_path,omitempty"`

	// PrefixedVars is a Hash from template variable names to output variable names.
	// It is used for values that must be unique across test runs (e.g., resource names).
	// In generated tests, the value will be prefixed (e.g., "tf-test-") and have a
	// random suffix appended. In documentation, the value is used verbatim.
	PrefixedVars map[string]string `yaml:"prefixed_vars,omitempty"`

	// Vars is a Hash from template variable names to output variable names.
	// It is used for values that should be inserted into the template as literal,
	// unmodified strings (e.g., labels, descriptions, or other non-identifier fields).
	// These values are used verbatim in both generated tests and documentation
	Vars map[string]string `yaml:"vars,omitempty"`

	// Some variables need to hold special values during tests, and cannot
	// be inferred by Open in Cloud Shell.  For instance, org_id
	// needs to be the correct value during integration tests, or else
	// org tests cannot pass. Other examples include an existing project_id,
	// a zone, a service account name, etc.
	//
	// test_env_vars is a Hash from template variable names to one of the
	// following symbols:
	//  - PROJECT_NAME
	//  - CREDENTIALS
	//  - REGION
	//  - ORG_ID
	//  - ORG_TARGET
	//  - BILLING_ACCT
	//  - MASTER_BILLING_ACCT
	//  - SERVICE_ACCT
	//  - CUST_ID
	//  - IDENTITY_USER
	//  - CHRONICLE_ID
	//  - VMWAREENGINE_PROJECT
	// This list corresponds to the `get*FromEnv` methods in provider_test.go.
	TestEnvVars map[string]string `yaml:"test_env_vars,omitempty"`

	// Hash to provider custom override values for generating test config
	// If field my-var is set in this hash, it will replace vars[my-var] in
	// tests. i.e. if vars["network"] = "my-vpc", without override:
	//   - doc config will have `network = "my-vpc"`
	//   - tests config will have `"network = my-vpc%{random_suffix}"`
	//     with context
	//       map[string]interface{}{
	//         "random_suffix": acctest.RandString()
	//       }
	//
	// If test_vars_overrides["network"] = "nameOfVpc()"
	//   - doc config will have `network = "my-vpc"`
	//   - tests will replace with `"network = %{network}"` with context
	//       map[string]interface{}{
	//         "network": nameOfVpc
	//         ...
	//       }
	TestVarsOverrides map[string]string `yaml:"test_vars_overrides,omitempty"`

	// Hash to provider custom override values for generating oics config
	// See test_vars_overrides for more details
	OicsVarsOverrides map[string]string `yaml:"oics_vars_overrides,omitempty"`

	// The version name of the test step's version if it's different than the
	// test version, eg. `beta`
	MinVersion string `yaml:"min_version,omitempty"`

	// Extra properties to ignore read on during import.
	// These properties will likely be custom code.
	IgnoreReadExtra []string `yaml:"ignore_read_extra,omitempty"`

	// Whether to skip import tests for this test step
	ExcludeImportTest bool `yaml:"exclude_import_test,omitempty"`

	// Whether to skip generating docs for this test step
	ExcludeDocs bool `yaml:"exclude_docs,omitempty"`

	DocumentationHCLText string `yaml:"-"`
	TestHCLText          string `yaml:"-"`
	OicsHCLText          string `yaml:"-"`
	PrimaryResourceId    string `yaml:"-"`
}

func (s *Step) TestStepSlug(productName, resourceName string) string {
	ret := fmt.Sprintf("%s%s_%sExample", productName, resourceName, google.Camelize(s.Name, "lower"))
	return ret
}

func (s *Step) Validate(rName, sName string) (es []error) {
	// TODO: Add check identifier when it's implemented
	if s.Name == "" {
		es = append(es, fmt.Errorf("missing `name` for one step in test sample %s in resource %s", sName, rName))
	}

	return es
}

func validateRegexForContents(r *regexp.Regexp, contents string, configPath string, objName string, vars map[string]string) {
	matches := r.FindAllStringSubmatch(contents, -1)
	for _, v := range matches {
		found := false
		for k, _ := range vars {
			if k == v[1] {
				found = true
				break
			}
		}
		if !found {
			log.Fatalf("Failed to find %s environment variable defined in YAML file when validating the file %s. Please define this in %s", v[1], configPath, objName)
		}
	}
}

// Executes step configuration templates for documentation and tests
func (s *Step) SetHCLText(sysfs fs.FS) {
	originalPrefixedVars := s.PrefixedVars
	// originalVars := s.Vars
	originalTestEnvVars := s.TestEnvVars
	docTestEnvVars := make(map[string]string)
	docs_defaults := map[string]string{
		"PROJECT_NAME":         "my-project-name",
		"PROJECT_NUMBER":       "1111111111111",
		"CREDENTIALS":          "my/credentials/filename.json",
		"REGION":               "us-west1",
		"ORG_ID":               "123456789",
		"ORG_DOMAIN":           "example.com",
		"ORG_TARGET":           "123456789",
		"BILLING_ACCT":         "000000-0000000-0000000-000000",
		"MASTER_BILLING_ACCT":  "000000-0000000-0000000-000000",
		"SERVICE_ACCT":         "my@service-account.com",
		"CUST_ID":              "A01b123xz",
		"IDENTITY_USER":        "cloud_identity_user",
		"PAP_DESCRIPTION":      "description",
		"CHRONICLE_ID":         "00000000-0000-0000-0000-000000000000",
		"VMWAREENGINE_PROJECT": "my-vmwareengine-project",
	}

	// Apply doc defaults to test_env_vars from YAML
	for key := range s.TestEnvVars {
		docTestEnvVars[key] = docs_defaults[s.TestEnvVars[key]]
	}
	s.TestEnvVars = docTestEnvVars
	s.DocumentationHCLText = s.ExecuteTemplate(sysfs)
	s.DocumentationHCLText = regexp.MustCompile(`\n\n$`).ReplaceAllString(s.DocumentationHCLText, "\n")

	// Remove region tags
	re1 := regexp.MustCompile(`# \[[a-zA-Z_ ]+\]\n`)
	re2 := regexp.MustCompile(`\n# \[[a-zA-Z_ ]+\]`)
	s.DocumentationHCLText = re1.ReplaceAllString(s.DocumentationHCLText, "")
	s.DocumentationHCLText = re2.ReplaceAllString(s.DocumentationHCLText, "")

	testPrefixedVars := make(map[string]string)
	testVars := make(map[string]string)
	testTestEnvVars := make(map[string]string)
	// Override prefixed_vars to inject test values into configs - will have
	//   - "a-example-var-value%{random_suffix}""
	//   - "%{my_var}" for overrides that have custom Golang values
	for key, value := range originalPrefixedVars {
		var newVal string
		if strings.Contains(value, "-") {
			newVal = fmt.Sprintf("tf-test-%s", value)
		} else if strings.Contains(value, "_") {
			newVal = fmt.Sprintf("tf_test_%s", value)
		} else {
			// Some vars like descriptions shouldn't have prefix
			newVal = value
		}
		// Random suffix is 10 characters and standard name length <= 64
		if len(newVal) > 54 {
			newVal = newVal[:54]
		}
		testPrefixedVars[key] = fmt.Sprintf("%s%%{random_suffix}", newVal)
	}

	// Apply overrides from YAML
	for key := range s.TestVarsOverrides {
		testPrefixedVars[key] = fmt.Sprintf("%%{%s}", key)
		testVars[key] = fmt.Sprintf("%%{%s}", key)
	}

	for key := range originalTestEnvVars {
		testTestEnvVars[key] = fmt.Sprintf("%%{%s}", key)
	}

	s.PrefixedVars = testPrefixedVars
	s.TestEnvVars = testTestEnvVars
	s.TestHCLText = s.ExecuteTemplate(sysfs)
	s.TestHCLText = regexp.MustCompile(`\n\n$`).ReplaceAllString(s.TestHCLText, "\n")
	// Remove region tags
	s.TestHCLText = re1.ReplaceAllString(s.TestHCLText, "")
	s.TestHCLText = re2.ReplaceAllString(s.TestHCLText, "")
	s.TestHCLText = SubstituteTestPaths(s.TestHCLText)

	// Reset the step
	s.PrefixedVars = originalPrefixedVars
	s.TestEnvVars = originalTestEnvVars
}

func (s *Step) ExecuteTemplate(sysfs fs.FS) string {
	templateContent, err := fs.ReadFile(sysfs, s.ConfigPath)
	if err != nil {
		glog.Exit(err)
	}

	fileContentString := string(templateContent)

	// Check that any variables in PrefixedVars, Vars or TestEnvVars used in the step are defined via YAML
	envVarRegex := regexp.MustCompile(`{{index \$\.TestEnvVars "([a-zA-Z_]*)"}}`)
	validateRegexForContents(envVarRegex, fileContentString, s.ConfigPath, "test_env_vars", s.TestEnvVars)
	varRegex := regexp.MustCompile(`{{index \$\.Vars "([a-zA-Z_]*)"}}`)
	validateRegexForContents(varRegex, fileContentString, s.ConfigPath, "vars", s.Vars)
	prefixedVarRegex := regexp.MustCompile(`{{index \$\.PrefixedVars "([a-zA-Z_]*)"}}`)
	validateRegexForContents(prefixedVarRegex, fileContentString, s.ConfigPath, "prefixed_vars", s.PrefixedVars)

	templateFileName := filepath.Base(s.ConfigPath)

	tmpl, err := template.New(templateFileName).Funcs(google.TemplateFunctions(sysfs)).Parse(fileContentString)
	if err != nil {
		glog.Exit(err)
	}

	contents := bytes.Buffer{}
	if err = tmpl.ExecuteTemplate(&contents, templateFileName, s); err != nil {
		glog.Exit(err)
	}

	rs := contents.String()

	if !strings.HasSuffix(rs, "\n") {
		rs = fmt.Sprintf("%s\n", rs)
	}

	return rs
}

func (s *Step) OiCSLink() string {
	v := url.Values{}
	v.Add("cloudshell_git_repo", "https://github.com/terraform-google-modules/docs-examples.git")
	v.Add("cloudshell_working_dir", s.Name)
	v.Add("cloudshell_image", "gcr.io/cloudshell-images/cloudshell:latest")
	v.Add("open_in_editor", "main.tf")
	v.Add("cloudshell_print", "./motd")
	v.Add("cloudshell_tutorial", "./tutorial.md")
	u := url.URL{
		Scheme:   "https",
		Host:     "console.cloud.google.com",
		Path:     "/cloudshell/open",
		RawQuery: v.Encode(),
	}
	return u.String()
}

func SubstituteExamplePaths(config string) string {
	config = strings.ReplaceAll(config, "../static/img/header-logo.png", "../static/header-logo.png")
	config = strings.ReplaceAll(config, "path/to/private.key", "../static/ssl_cert/test.key")
	config = strings.ReplaceAll(config, "path/to/id_rsa.pub", "../static/ssh_rsa.pub")
	config = strings.ReplaceAll(config, "path/to/certificate.crt", "../static/ssl_cert/test.crt")
	return config
}

func SubstituteTestPaths(config string) string {
	config = strings.ReplaceAll(config, "../static/img/header-logo.png", "test-fixtures/header-logo.png")
	config = strings.ReplaceAll(config, "path/to/private.key", "test-fixtures/test.key")
	config = strings.ReplaceAll(config, "path/to/certificate.crt", "test-fixtures/test.crt")
	config = strings.ReplaceAll(config, "path/to/index.zip", "%{zip_path}")
	config = strings.ReplaceAll(config, "verified-domain.com", "tf-test-domain%{random_suffix}.gcp.tfacc.hashicorptest.com")
	config = strings.ReplaceAll(config, "path/to/id_rsa.pub", "test-fixtures/ssh_rsa.pub")
	return config
}

// Executes step configuration templates for documentation and tests
func (s *Step) SetOiCSHCLText(sysfs fs.FS) {
	originalPrefixedVars := s.PrefixedVars

	// // Remove region tags
	re1 := regexp.MustCompile(`# \[[a-zA-Z_ ]+\]\n`)
	re2 := regexp.MustCompile(`\n# \[[a-zA-Z_ ]+\]`)

	testPrefixedVars := make(map[string]string)
	for key, value := range originalPrefixedVars {
		testPrefixedVars[key] = fmt.Sprintf("%s-${local.name_suffix}", value)
	}

	// Apply overrides from YAML
	for key, value := range s.OicsVarsOverrides {
		testPrefixedVars[key] = value
	}

	s.PrefixedVars = testPrefixedVars
	s.OicsHCLText = s.ExecuteTemplate(sysfs)
	s.OicsHCLText = regexp.MustCompile(`\n\n$`).ReplaceAllString(s.OicsHCLText, "\n")

	// Remove region tags
	s.OicsHCLText = re1.ReplaceAllString(s.OicsHCLText, "")
	s.OicsHCLText = re2.ReplaceAllString(s.OicsHCLText, "")
	s.OicsHCLText = SubstituteExamplePaths(s.OicsHCLText)

	// Reset the step
	s.PrefixedVars = originalPrefixedVars
}
