// Copyright 2024 Google Inc.
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
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
	"github.com/golang/glog"
	"gopkg.in/yaml.v3"
)

// Generates configs to be shown as examples in docs and outputted as tests
// from a shared template
type Examples struct {
	// google.YamlValidator

	// include Compile::Core
	// include Google::GolangUtils

	// The name of the example in lower snake_case.
	// Generally takes the form of the resource name followed by some detail
	// about the specific test. For example, "address_with_subnetwork".
	Name string

	// The id of the "primary" resource in an example. Used in import tests.
	// This is the value that will appear in the Terraform config url. For
	// example:
	// resource "google_compute_address" {{primary_resource_id}} {
	//   ...
	// }
	PrimaryResourceId string `yaml:"primary_resource_id"`

	// Optional resource type of the "primary" resource. Used in import tests.
	// If set, this will override the default resource type implied from the
	// object parent
	PrimaryResourceType string `yaml:"primary_resource_type"`

	// Vars is a Hash from template variable names to output variable names.
	// It will use the provided value as a prefix for generated tests, and
	// insert it into the docs verbatim.
	Vars map[string]string

	// Some variables need to hold special values during tests, and cannot
	// be inferred by Open in Cloud Shell.  For instance, org_id
	// needs to be the correct value during integration tests, or else
	// org tests cannot pass. Other examples include an existing project_id,
	// a zone, a service account name, etc.
	//
	// test_env_vars is a Hash from template variable names to one of the
	// following symbols:
	//  - :PROJECT_NAME
	//  - :CREDENTIALS
	//  - :REGION
	//  - :ORG_ID
	//  - :ORG_TARGET
	//  - :BILLING_ACCT
	//  - :MASTER_BILLING_ACCT
	//  - :SERVICE_ACCT
	//  - :CUST_ID
	//  - :IDENTITY_USER
	// This list corresponds to the `get*FromEnv` methods in provider_test.go.
	TestEnvVars map[string]string `yaml:"test_env_vars"`

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
	TestVarsOverrides map[string]string `yaml:"test_vars_overrides"`

	// Hash to provider custom override values for generating oics config
	// See test_vars_overrides for more details
	OicsVarsOverrides map[string]string `yaml:"oics_vars_overrides"`

	// The version name of of the example's version if it's different than the
	// resource version, eg. `beta`
	//
	// This should be the highest version of all the features used in the
	// example; if there's a single beta field in an example, the example's
	// min_version is beta. This is only needed if an example uses features
	// with a different version than the resource; a beta resource's examples
	// are all automatically versioned at beta.
	//
	// When an example has a version of beta, each resource must use the
	// `google-beta` provider in the config. If the `google` provider is
	// implicitly used, the test will fail.
	//
	// NOTE: Until Terraform 0.12 is released and is used in the OiCS tests, an
	// explicit provider block should be defined. While the tests @ 0.12 will
	// use `google-beta` automatically, past Terraform versions required an
	// explicit block.
	MinVersion string `yaml:"min_version"`

	// Extra properties to ignore read on during import.
	// These properties will likely be custom code.
	IgnoreReadExtra []string `yaml:"ignore_read_extra"`

	// Whether to skip generating tests for this resource
	SkipTest bool `yaml:"skip_test"`

	// Whether to skip generating docs for this example
	SkipDocs bool `yaml:"skip_docs"`

	// Whether to skip import tests for this example
	SkipImportTest bool `yaml:"skip_import_test"`

	// The name of the primary resource for use in IAM tests. IAM tests need
	// a reference to the primary resource to create IAM policies for
	PrimaryResourceName string `yaml:"primary_resource_name"`

	// The name of the location/region override for use in IAM tests. IAM
	// tests may need this if the location is not inherited on the resource
	// for one reason or another
	RegionOverride string `yaml:"region_override"`

	// The path to this example's Terraform config.
	// Defaults to `templates/terraform/examples/{{name}}.tf.erb`
	ConfigPath string `yaml:"config_path"`

	// If the example should be skipped during VCR testing.
	// This is the case when something about the resource or config causes VCR to fail for example
	// a resource with a unique identifier generated within the resource via id.UniqueId()
	// Or a config with two fine grained resources that have a race condition during create
	SkipVcr bool `yaml:"skip_vcr"`

	// Specify which external providers are needed for the testcase.
	// Think before adding as there is latency and adds an external dependency to
	// your test so avoid if you can.
	ExternalProviders []string `yaml:"external_providers"`

	DocumentationHCLText string
	TestHCLText          string
}

// Set default value for fields
func (e *Examples) UnmarshalYAML(n *yaml.Node) error {
	type exampleAlias Examples
	aliasObj := (*exampleAlias)(e)

	err := n.Decode(&aliasObj)
	if err != nil {
		return err
	}

	if e.ConfigPath == "" {
		e.ConfigPath = fmt.Sprintf("templates/terraform/examples/go/%s.tf.tmpl", e.Name)
	}
	e.SetHCLText()

	return nil
}

// Executes example templates for documentation and tests
func (e *Examples) SetHCLText() {
	originalVars := e.Vars
	originalTestEnvVars := e.TestEnvVars
	docTestEnvVars := make(map[string]string)
	docs_defaults := map[string]string{
		"PROJECT_NAME":        "my-project-name",
		"CREDENTIALS":         "my/credentials/filename.json",
		"REGION":              "us-west1",
		"ORG_ID":              "123456789",
		"ORG_DOMAIN":          "example.com",
		"ORG_TARGET":          "123456789",
		"BILLING_ACCT":        "000000-0000000-0000000-000000",
		"MASTER_BILLING_ACCT": "000000-0000000-0000000-000000",
		"SERVICE_ACCT":        "my@service-account.com",
		"CUST_ID":             "A01b123xz",
		"IDENTITY_USER":       "cloud_identity_user",
		"PAP_DESCRIPTION":     "description",
	}

	// Apply doc defaults to test_env_vars from YAML
	for key := range e.TestEnvVars {
		docTestEnvVars[key] = docs_defaults[e.TestEnvVars[key]]
	}
	e.TestEnvVars = docTestEnvVars
	e.DocumentationHCLText = ExecuteTemplate(e, e.ConfigPath, true)
	e.DocumentationHCLText = regexp.MustCompile(`\n\n$`).ReplaceAllString(e.DocumentationHCLText, "\n")

	// Remove region tags
	re1 := regexp.MustCompile(`# \[[a-zA-Z_ ]+\]\n`)
	re2 := regexp.MustCompile(`\n# \[[a-zA-Z_ ]+\]`)
	e.DocumentationHCLText = re1.ReplaceAllString(e.DocumentationHCLText, "")
	e.DocumentationHCLText = re2.ReplaceAllString(e.DocumentationHCLText, "")

	testVars := make(map[string]string)
	testTestEnvVars := make(map[string]string)
	// Override vars to inject test values into configs - will have
	//   - "a-example-var-value%{random_suffix}""
	//   - "%{my_var}" for overrides that have custom Golang values
	for key, value := range originalVars {
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
		testVars[key] = fmt.Sprintf("%s%%{random_suffix}", newVal)
	}

	// Apply overrides from YAML
	for key := range e.TestVarsOverrides {
		testVars[key] = fmt.Sprintf("%%{%s}", key)
	}
	for key := range originalTestEnvVars {
		testTestEnvVars[key] = fmt.Sprintf("%%{%s}", key)
	}

	e.Vars = testVars
	e.TestEnvVars = testTestEnvVars
	e.TestHCLText = ExecuteTemplate(e, e.ConfigPath, true)
	e.TestHCLText = regexp.MustCompile(`\n\n$`).ReplaceAllString(e.TestHCLText, "\n")
	// Remove region tags
	e.TestHCLText = re1.ReplaceAllString(e.TestHCLText, "")
	e.TestHCLText = re2.ReplaceAllString(e.TestHCLText, "")
	e.TestHCLText = SubstituteTestPaths(e.TestHCLText)

	// Reset the example
	e.Vars = originalVars
	e.TestEnvVars = originalTestEnvVars
}

func ExecuteTemplate(e any, templatePath string, appendNewline bool) string {
	templates := []string{
		templatePath,
		"templates/terraform/expand_resource_ref.tmpl",
		"templates/terraform/custom_flatten/go/bigquery_table_ref.go.tmpl",
		"templates/terraform/flatten_property_method.go.tmpl",
		"templates/terraform/expand_property_method.go.tmpl",
		"templates/terraform/update_mask.go.tmpl",
		"templates/terraform/nested_query.go.tmpl",
		"templates/terraform/unordered_list_customize_diff.go.tmpl",
	}
	templateFileName := filepath.Base(templatePath)

	tmpl, err := template.New(templateFileName).Funcs(google.TemplateFunctions).ParseFiles(templates...)
	if err != nil {
		glog.Exit(err)
	}

	contents := bytes.Buffer{}
	if err = tmpl.ExecuteTemplate(&contents, templateFileName, e); err != nil {
		glog.Exit(err)
	}

	rs := contents.String()

	if !strings.HasSuffix(rs, "\n") && appendNewline {
		rs = fmt.Sprintf("%s\n", rs)
	}

	return rs
}

func (e *Examples) OiCSLink() string {
	v := url.Values{}
	// TODO Q2: Values.Encode() sorts the values by key alphabetically. This will produce
	//			diffs for every URL when we convert to using this function. We should sort the
	// 			Ruby-version query alphabetically beforehand to remove these diffs.
	v.Add("cloudshell_git_repo", "https://github.com/terraform-google-modules/docs-examples.git")
	v.Add("cloudshell_working_dir", e.Name)
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

func (e *Examples) TestSlug(productName, resourceName string) string {
	ret := fmt.Sprintf("%s%s_%sExample", productName, resourceName, google.Camelize(e.Name, "lower"))
	return ret
}

func (e *Examples) ResourceType(terraformName string) string {
	if e.PrimaryResourceType != "" {
		return e.PrimaryResourceType
	}
	return terraformName
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

// func (e *Examples) validate() {
// super
// check :name, type: String, required: true
// check :primary_resource_id, type: String
// check :min_version, type: String
// check :vars, type: Hash
// check :test_env_vars, type: Hash
// check :test_vars_overrides, type: Hash
// check :ignore_read_extra, type: Array, item_type: String, default: []
// check :primary_resource_name, type: String
// check :skip_test, type: TrueClass
// check :skip_import_test, type: TrueClass
// check :skip_docs, type: TrueClass
// check :config_path, type: String, default: "templates/terraform/examples///{name}.tf.erb"
// check :skip_vcr, type: TrueClass
// }

// TODO
// validate_external_providers

// func (e *Examples) merge(other) {
// result = self.class.new
// instance_variables.each do |v|
//   result.instance_variable_set(v, instance_variable_get(v))
// end

// other.instance_variables.each do |v|
//   if other.instance_variable_get(v).instance_of?(Array)
//     result.instance_variable_set(v, deep_merge(result.instance_variable_get(v),
//                                                other.instance_variable_get(v)))
//   else
//     result.instance_variable_set(v, other.instance_variable_get(v))
//   end
// end

// result
// }
